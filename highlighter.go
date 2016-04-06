package color

import (
	"sync"
	"unicode"
)

// stateFn represents the state of the highlighter as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// highlighter holds the state of the scanner.
type highlighter struct {
	s     string // string being scanned
	buf   []byte
	pos   int    // position in buf
	last  int    // position after last verb
	attrs []byte // attributes of current highlight verb
}

var hlPool = sync.Pool{
	New: func() interface{} { return new(highlighter) },
}

// Highlight replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string.
func Highlight(s string) string {
	h := get(s)
	h.run()
	return h.free()
}

func get(s string) (hl *highlighter) {
	hl = hlPool.Get().(*highlighter)
	hl.s = s
	return
}

func (hl *highlighter) free() string {
	s := string(hl.buf)
	hl.buf = hl.buf[:0]
	hl.pos = 0
	hlPool.Put(hl)
	return s
}

// run runs the state machine for the highlighter.
func (hl *highlighter) run() {
	for state := scanText; state != nil; {
		state = state(hl)
	}
}

// get returns the current rune.
func (hl *highlighter) get() rune {
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
	hl.last = hl.pos
	for ; hl.pos < len(hl.s); hl.pos++ {
		if hl.get() != '%' {
			continue
		}
		if hl.last < hl.pos {
			hl.buf = append(hl.buf, hl.s[hl.last:hl.pos]...)
		}
		hl.pos++
		if hl.pos >= len(hl.s) {
			return nil
		}
		switch hl.get() {
		case 'r':
			hl.pos++
			return verbReset
		case 'h':
			hl.pos += 2
			return scanHighlight
		}
	}
	return nil
}

// verbReset appen the reset verb with the reset control sequence.
func verbReset(hl *highlighter) stateFn {
	hl.attrs = append(hl.attrs, attrs["reset"]...)
	hl.appendAttrs()
	hl.attrs = hl.attrs[:0]
	return scanText
}

// scanHighlight scans the highlight verb for attributes,
// then replaces it with a control sequence derived from said attributes.
func scanHighlight(hl *highlighter) stateFn {
	for ; hl.pos < len(hl.s); hl.pos++ {
		r := hl.get()
		switch {
		case r == 'f':
			return scanColor256(hl, preFg256)
		case r == 'b':
			return scanColor256(hl, preBg256)
		case unicode.IsLetter(r):
			return scanAttribute(hl, 0)
		case r == '+':
			// skip
		case r == ']':
			if len(hl.attrs) != 0 {
				hl.appendAttrs()
			}
			hl.pos++
			fallthrough
		default:
			// reuse this buffer
			hl.attrs = hl.attrs[:0]
			return scanText
		}
	}
	return nil
}

// scanAttribute scans a named attribute
func scanAttribute(hl *highlighter, off int) stateFn {
	start := hl.pos - off
	for ; hl.pos < len(hl.s); hl.pos++ {
		if !unicode.IsLetter(hl.get()) {
			if a, ok := attrs[hl.s[start:hl.pos]]; ok {
				hl.attrs = append(hl.attrs, a...)
			}
			return scanHighlight
		}
	}
	return nil
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
	for ; hl.pos < len(hl.s); hl.pos++ {
		if !unicode.IsNumber(hl.get()) {
			hl.attrs = append(hl.attrs, pre...)
			hl.attrs = append(hl.attrs, hl.s[start:hl.pos]...)
			return scanHighlight
		}
	}
	return nil
}
