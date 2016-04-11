package color

import (
	"strings"
	"sync"
	"unicode"
)

const (
	errMissing = "%%!h(MISSING)" // no attributes in the highlight verb
	errInvalid = "%%!h(INVALID)" // invalid character in the highlight verb
	errBadAttr = "%%!h(BADATTR)" // unknown attribute in the highlight verb
)

// highlighter holds the state of the scanner.
type highlighter struct {
	s     string // string being scanned
	pos   int    // position in s
	buf   buffer // where result is built
	attrs buffer // attributes of current verb
}

// Shighlightf replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string.
// This is a low level function that only handles highlight verbs, you should
// use color.Sprintf most of the time as it wraps around fmt.Sprintf which
// handles other verbs.
func Shighlightf(s string) string {
	hl := getHighlighter(s)
	defer hl.free()
	hl.run(initHighlight)
	s = string(hl.buf)
	return s
}

func Sstripf(s string) string {
	hl := getHighlighter(s)
	defer hl.free()
	hl.run(initStrip)
	s = string(hl.buf)
	return s
}

// highlighterPool allows the reuse of highlighters to avoid allocations.
var highlighterPool = sync.Pool{
	New: func() interface{} {
		hl := new(highlighter)
		// The initial capacities avoid constant reallocation during growth.
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

// free resets the highlighter.
func (hl *highlighter) free() {
	hl.buf.reset()
	hl.pos = 0
	highlighterPool.Put(hl)
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// run runs the state machine for the highlighter.
func (hl *highlighter) run(init stateFn) {
	for state := init; state != nil; {
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

// writePrev writes n previous characters to the buffer
func (hl *highlighter) writePrev(n int) {
	hl.pos++
	hl.buf.writeString(hl.s[hl.pos-n : hl.pos])
}

func initHighlight(hl *highlighter) stateFn {
	return scanText(hl, scanVerb)
}

func initStrip(hl *highlighter) stateFn {
	return scanText(hl, stripVerb)
}

// scanText scans until the next verb.
func scanText(hl *highlighter, fn stateFn) stateFn {
	ppos := hl.pos
	// Find next verb.
	for {
		switch hl.get() {
		case eof:
			if hl.pos > ppos {
				// Append remaining characters.
				hl.buf.writeString(hl.s[ppos:hl.pos])
			}
			return nil
		case '%':
			if hl.pos > ppos {
				// Append the characters after the last verb.
				hl.buf.writeString(hl.s[ppos:hl.pos])
			}
			hl.pos++
			return fn

		}
		hl.pos++
	}
}

func stripVerb(hl *highlighter) stateFn {
	switch hl.get() {
	case 'r':
		// Strip the reset verb.
		hl.pos++
	case 'h':
		// Strip inside the highlight verb.
		hl.pos++
		j := strings.IndexByte(hl.s[hl.pos:], ']')
		if j == -1 {
			hl.buf.writeString(errInvalid)
			return nil
		}
		hl.pos += j + 1
	case eof:
		// Let fmt handle "%!h(NOVERB)".
		hl.buf.writeByte('%')
		return nil
	default:
		// Include the verb.
		hl.writePrev(2)
	}
	return initStrip
}

// scanVerb scans the current verb.
func scanVerb(hl *highlighter) stateFn {
	switch hl.get() {
	case 'r':
		hl.pos++
		return verbReset
	case 'h':
		hl.pos += 2
		return scanHighlight
	case eof:
		// Let fmt handle "%!h(NOVERB)".
		hl.buf.writeByte('%')
		return nil
	}
	hl.writePrev(2)
	return initHighlight
}

// verbReset writes the reset verb with the reset control sequence.
func verbReset(hl *highlighter) stateFn {
	hl.attrs.writeString(attrs["reset"])
	hl.writeAttrs()
	hl.attrs.reset()
	return initHighlight
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
		return initHighlight
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
			return initHighlight
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
