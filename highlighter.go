package color

import (
	"sync"
	"unicode"
)

const (
	errInvalid = "%!h(INVALID)" // something unexpected
	errMissing = "%!h(MISSING)" // no attrs
	errBadAttr = "%!h(BADATTR)" // attr isn't a color or in the map
	errNoVerb  = "%!(NOVERB)"   // no verb
)

// stateFn represents the state of the highlighter as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// highlighter holds the state of the scanner.
type highlighter struct {
	s     string // string being scanned
	buf   buffer // buffer for result
	pos   int    // position in buf
	attrs buffer // attributes of current highlight verb
}

var hlPool = sync.Pool{
	New: func() interface{} {
		hl := new(highlighter)
		hl.buf = make([]byte, 0, 30)
		hl.attrs = make([]byte, 0, 10)
		return hl
	},
}

// Highlight replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string.
func Highlight(s string) string {
	hl := getHighlighter(s)
	hl.run()
	return hl.free()
}

func getHighlighter(s string) (hl *highlighter) {
	hl = hlPool.Get().(*highlighter)
	hl.s = s
	return
}

func (hl *highlighter) free() (s string) {
	s = string(hl.buf)
	hl.buf.reset()
	hl.pos = 0
	hlPool.Put(hl)
	return
}

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

// appends a control sequence derived from h.attrs[1:] to h.buf.
func (hl *highlighter) writeAttrs() {
	hl.buf.writeString(csi)
	hl.buf.write(hl.attrs[1:])
	hl.buf.writeByte('m')
}

func (hl *highlighter) writePrev(n int) {
	hl.buf.writeString(hl.s[n:hl.pos])
}

// scanText scans until the next highlight or reset verb.
func scanText(hl *highlighter) stateFn {
	// previous position
	ppos := hl.pos
	for {
		if r := hl.get(); r == eof {
			if ppos < hl.pos {
				hl.writePrev(ppos)
			}
			return nil
		} else if r == '%' {
			if ppos < hl.pos {
				hl.writePrev(ppos)
			}
			break
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
		hl.buf.writeString(errNoVerb)
		return nil
	}
	hl.pos++
	hl.writePrev(hl.pos - 2)
	return scanText
}

// verbReset appends the reset verb with the reset control sequence.
func verbReset(hl *highlighter) stateFn {
	hl.attrs.writeString(attrs["reset"])
	hl.writeAttrs()
	hl.attrs.reset()
	return scanText
}

// scanHighlight scans the highlight verb for attributes,
// then replaces it with a control sequence derived from said attributes.
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

// scanAttribute scans a named attribute
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

// scanColor256 scans a 256 color attribute
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
