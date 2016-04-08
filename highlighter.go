package color

import (
	"strings"
	"sync"
	"unicode"
)

// see doc.go for an explanation of these
const (
	errInvalid = "%%!h(INVALID)"
	errMissing = "%%!h(MISSING)"
	errBadAttr = "%%!h(BADATTR)"
)

// highlighter holds the state of the scanner.
type highlighter struct {
	s     string // string being scanned
	pos   int    // position in s
	buf   buffer // result
	attrs buffer // attributes of current verb
}

// Highlight replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string.
// This is a low-level function that only scans highlight verbs. The color.Printf functions
// are the intended user functions as they wrap around the fmt.Printf functions,
// which handle the rest. Only use this for performance reasons.
func Highlight(s string) string {
	hl := getHighlighter(s)
	hl.run()
	return string(hl.free())
}

// highlighterPool reuses highlighter objects to avoid an allocation per invocation.
var highlighterPool = sync.Pool{
	New: func() interface{} {
		hl := new(highlighter)
		// initial capacities avoid constant reallocation during growth.
		hl.buf = make([]byte, 0, 30)
		hl.attrs = make([]byte, 0, 10)
		return hl
	},
}

// getHighlighter returns a new initialized highlighter from the pool.
func getHighlighter(s string) (hl *highlighter) {
	hl = highlighterPool.Get().(*highlighter)
	hl.s = s
	return
}

// free resets the highlighter and returns the buffer.
func (hl *highlighter) free() (b []byte) {
	b = hl.buf
	hl.buf.reset()
	hl.pos = 0
	highlighterPool.Put(hl)
	return
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// run runs the state machine for the highlighter.
func (hl *highlighter) run() {
	for state := scanText; state != nil; {
		state = state(hl)
	}
}

const eof = -1

// get returns the current rune.
func (hl *highlighter) get() rune {
	if hl.pos >= len(hl.s) {
		return eof
	}
	return rune(hl.s[hl.pos])
}

// writeAttrs writes a control sequence derived from h.attrs[1:] to h.buf.
func (hl *highlighter) writeAttrs() {
	hl.buf.writeString(csi)
	hl.buf.write(hl.attrs[1:])
	hl.buf.writeByte('m')
}

// writePrev writes n previous characters to the buffer.
func (hl *highlighter) writePrev(n int) {
	hl.buf.writeString(hl.s[n:hl.pos])
}

// scanText scans until the next highlight or reset verb.
func scanText(hl *highlighter) stateFn {
	// previous position
	ppos := hl.pos
LOOP:
	for {
		switch hl.get() {
		case eof:
			if hl.pos > ppos {
				hl.writePrev(ppos)
			}
			return nil
		case '%':
			if hl.pos > ppos {
				hl.writePrev(ppos)
			}
			break LOOP

		}
		hl.pos++
	}
	hl.pos++
	switch hl.get() {
	case 'r':
		hl.pos++
		return verbReset
	case 'h':
		hl.pos += 2
		return scanHighlight
	case eof:
		// let fmt handle NOVERB
		hl.buf.writeByte('%')
		return nil
	}
	hl.pos++
	hl.writePrev(hl.pos - 2)
	return scanText
}

// verbReset writes the reset verb with the reset control sequence.
func verbReset(hl *highlighter) stateFn {
	hl.attrs.writeString(attrs["reset"])
	hl.writeAttrs()
	hl.attrs.reset()
	return scanText
}

// scanHighlight scans the highlight verb for attributes,
// then writes a control sequence derived from said attributes to the buffer.
func scanHighlight(hl *highlighter) stateFn {
	r := hl.get()
	switch {
	case r == 'f':
		return scanColor256(hl, preFg256)
	case r == 'b':
		return scanColor256(hl, preBg256)
	case unicode.IsLetter(r):
		return scanAttribute(hl, 0)
	case r == '+':
		hl.pos++
		return scanHighlight
	case r == ']':
		if len(hl.attrs) != 0 {
			hl.writeAttrs()
		} else {
			hl.buf.writeString(errMissing)
		}
		hl.attrs.reset()
		hl.pos++
		return scanText
	default:
		return abortHighlight(hl, errInvalid)
	}
}

// scanAttribute scans a named attribute.
func scanAttribute(hl *highlighter, off int) stateFn {
	start := hl.pos - off
	for unicode.IsLetter(hl.get()) {
		hl.pos++
	}
	if a, ok := attrs[hl.s[start:hl.pos]]; ok {
		hl.attrs.writeString(a)
	} else {
		return abortHighlight(hl, errBadAttr)
	}
	return scanHighlight
}

// abortHighlight writes a error to the buffer and
// then skips to the end of the highlight verb.
func abortHighlight(hl *highlighter, msg string) stateFn {
	hl.buf.writeString(msg)
	hl.attrs.reset()
	for {
		switch hl.get() {
		case ']':
			hl.pos++
			return scanText
		case eof:
			return nil
		}
		hl.pos++
	}
}

// scanColor256 scans a 256 color attribute.
func scanColor256(hl *highlighter, pre string) stateFn {
	hl.pos++
	if hl.get() != 'g' {
		return scanAttribute(hl, 1)
	}
	hl.pos++
	if !unicode.IsNumber(hl.get()) {
		return scanAttribute(hl, 2)
	}
	start := hl.pos
	hl.pos++
	for unicode.IsNumber(hl.get()) {
		hl.pos++
	}
	hl.attrs.writeString(pre)
	hl.attrs.writeString(hl.s[start:hl.pos])
	return scanHighlight
}

// bufferPool reuses buffers to avoid an allocation per invocation.
var bufferPool = sync.Pool{
	New: func() interface{} {
		// initial capacity avoids constant reallocation during growth.
		return buffer(make([]byte, 0, 30))
	},
}

// stripVerbs removes all highlight verbs in s.
func stripVerbs(s string) string {
	buf := bufferPool.Get().(buffer)
	// pi is the index after last verb.
	var pi, i int
LOOP:
	for ; ; i++ {
		if i >= len(s) {
			if i > pi {
				buf.writeString(s[pi:i])
			}
			break
		} else if s[i] != '%' {
			continue
		}
		if i > pi {
			buf.writeString(s[pi:i])
		}
		i++
		if i >= len(s) {
			// let fmt handle NOVERB
			buf.writeByte('%')
			break
		}
		switch s[i] {
		case 'r':
			// strip the reset verb
			pi = i + 1
		case 'h':
			// strip inside the highlight verb
			j := strings.IndexByte(s[i+1:], ']')
			if j == -1 {
				buf.writeString(errInvalid)
				break LOOP
			}
			i += j + 1
			pi = i + 1
		default:
			// include the verb
			pi = i - 1
		}
	}
	s = string(buf)
	buf.reset()
	bufferPool.Put(buf)
	return s
}
