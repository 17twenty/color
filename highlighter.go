package color

import (
	"sync"
	"unicode"
)

const (
	commaSpaceBytes  = ", "
	nilAngleBytes    = "<nil>"
	nilParenBytes    = "(nil)"
	nilBytes         = "nil"
	mapBytes         = "map["
	percentBangBytes = "%!"
	missingBytes     = "(MISSING)"
	badIndexBytes    = "(BADINDEX)"
	panicBytes       = "(PANIC="
	extraBytes       = "%!(EXTRA "
	irparenBytes     = "i)"
	bytesBytes       = "[]byte{"
	badWidthBytes    = "%!(BADWIDTH)"
	badPrecBytes     = "%!(BADPREC)"
	noVerbBytes      = "%!(NOVERB)"
)

// stateFn represents the state of the highlighter as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// highlighter holds the state of the scanner.
type highlighter struct {
	s     string // string being scanned
	buf   []byte // buffer for result
	pos   int    // position in buf
	attrs []byte // attributes of current highlight verb
}

var hlPool = sync.Pool{
	New: func() interface{} { return new(highlighter) },
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
	hl.buf = hl.buf[:0]
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
func (hl *highlighter) appendAttrs() {
	hl.buf = append(hl.buf, csi...)
	hl.buf = append(hl.buf, hl.attrs[1:]...)
	hl.buf = append(hl.buf, 'm')
}

// scanText scans until the next highlight or reset verb.
func scanText(hl *highlighter) stateFn {
	last := hl.pos
	for {
		if r := hl.get(); r == eof {
			return nil
		} else if r == '%' {
			break
		}
		hl.pos++
	}
	if last < hl.pos {
		hl.buf = append(hl.buf, hl.s[last:hl.pos]...)
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
		hl.buf = append(hl.buf, noVerbBytes...)
		return nil
	}
	hl.pos++
	return scanText
}

// verbReset appends the reset verb with the reset control sequence.
func verbReset(hl *highlighter) stateFn {
	hl.attrs = append(hl.attrs, attrs["reset"]...)
	hl.appendAttrs()
	hl.attrs = hl.attrs[:0]
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
			hl.appendAttrs()
		}
		hl.pos++
		fallthrough
	default:
		hl.attrs = hl.attrs[:0]
		return scanText
	}
}

// scanAttribute scans a named attribute
func scanAttribute(hl *highlighter, off int) stateFn {
	start := hl.pos - off
	for unicode.IsLetter(hl.get()) {
		hl.pos++
	}
	if a, ok := attrs[hl.s[start:hl.pos]]; ok {
		hl.attrs = append(hl.attrs, a...)
	}
	return scanHighlight
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
	hl.attrs = append(hl.attrs, pre...)
	hl.attrs = append(hl.attrs, hl.s[start:hl.pos]...)
	return scanHighlight
}
