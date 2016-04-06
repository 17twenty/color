package color

import "unicode"

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

// Highlight replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string.
func Highlight(s string) string {
	h := &highlighter{s: s}
	h.run()
	return string(h.buf)
}

// run runs the state machine for the highlighter.
func (h *highlighter) run() {
	for state := scanText; state != nil; {
		state = state(h)
	}
}

// get returns the current rune.
func (h *highlighter) get() rune {
	return rune(h.s[h.pos])
}

// appends a control sequence derived from h.attrs[1:] to h.buf.
func (h *highlighter) appendAttrs() {
	h.buf = append(h.buf, csi...)
	h.buf = append(h.buf, h.attrs[1:]...)
	h.buf = append(h.buf, 'm')
}

// scanText scans until the next highlight or reset verb.
func scanText(h *highlighter) stateFn {
	h.last = h.pos
	for ; h.pos < len(h.s); h.pos++ {
		if h.get() != '%' {
			continue
		}
		if h.last < h.pos {
			h.buf = append(h.buf, h.s[h.last:h.pos]...)
		}
		h.pos++
		if h.pos >= len(h.s) {
			return nil
		}
		switch h.get() {
		case 'r':
			h.pos++
			return verbReset
		case 'h':
			h.pos += 2
			return scanHighlight
		}
	}
	return nil
}

// verbReset appen the reset verb with the reset control sequence.
func verbReset(h *highlighter) stateFn {
	h.attrs = append(h.attrs, attrs["reset"]...)
	h.appendAttrs()
	return scanText
}

// scanHighlight scans the highlight verb for attributes,
// then replaces it with a control sequence derived from said attributes.
func scanHighlight(h *highlighter) stateFn {
	for ; h.pos < len(h.s); h.pos++ {
		r := h.get()
		switch {
		case r == 'f':
			return scanColor256(h, preFg256)
		case r == 'b':
			return scanColor256(h, preBg256)
		case unicode.IsLetter(r):
			return scanAttribute(h, 0)
		case r == '+':
			// skip
		case r == ']':
			if len(h.attrs) != 0 {
				h.appendAttrs()
			}
			h.pos++
			fallthrough
		default:
			// reuse this buffer
			h.attrs = h.attrs[0:0]
			return scanText
		}
	}
	return nil
}

// scanAttribute scans a named attribute
func scanAttribute(h *highlighter, off int) stateFn {
	start := h.pos - off
	for ; h.pos < len(h.s); h.pos++ {
		if !unicode.IsLetter(h.get()) {
			if a, ok := attrs[h.s[start:h.pos]]; ok {
				h.attrs = append(h.attrs, a...)
			}
			return scanHighlight
		}
	}
	return nil
}

// scanColor256 scans a 256 color attribute
func scanColor256(h *highlighter, pre string) stateFn {
	h.pos++
	if h.get() != 'g' {
		return scanAttribute(h, 1)
	}
	h.pos++
	if !unicode.IsNumber(h.get()) {
		return scanAttribute(h, 2)
	}
	start := h.pos
	for ; h.pos < len(h.s); h.pos++ {
		if !unicode.IsNumber(h.get()) {
			h.attrs = append(h.attrs, pre...)
			h.attrs = append(h.attrs, h.s[start:h.pos]...)
			return scanHighlight
		}
	}
	return nil
}
