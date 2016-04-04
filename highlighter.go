package color

import (
	"unicode"
	"unicode/utf8"
)

const eof = -1

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// highlighter holds the state of the scanner.
type highlighter struct {
	b     []byte // string being scanned
	pos   int    // position in s
	width int    // width of last rune read from s
	start int    // start position of current verb
	attrs []byte // attributes of current highlight verb
}

// Highlight replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string
func Highlight(s string) string {
	h := &highlighter{b: []byte(s)}
	h.run()
	return string(h.b)
}

// run runs the state machine for the highlighter.
func (h *highlighter) run() {
	for state := scanText; state != nil; {
		state = state(h)
	}
}

// next returns the next rune in the input.
func (h *highlighter) next() rune {
	if h.pos >= len(h.b) {
		h.width = 0
		return eof
	}
	r, w := utf8.DecodeRune(h.b[h.pos:])
	h.pos += w
	h.width = w
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (h *highlighter) backup() {
	h.pos -= h.width
}

// replaces the verb with a control sequence derived from h.attrs[1:].
func (h *highlighter) replace() {
	back := h.pos - h.start
	h.attrs = h.attrs[1:]
	b := make([]byte, len(h.b)-back+len(csi)+len(h.attrs)+1) // one more for the 'm'
	l := h.pos - back
	copy(b, h.b[:l])
	copy(b[l:], csi)
	l += len(csi)
	copy(b[l:], h.attrs)
	l += len(h.attrs)
	h.attrs = nil
	b[l] = 'm'
	l++
	copy(b[l:], h.b[h.pos:])
	h.pos = l
	h.b = b
}

// scans until the next highlight or reset verb.
func scanText(h *highlighter) stateFn {
	for {
		switch h.next() {
		case eof:
			return nil
		default:
			continue
		case '%':
		}
		switch h.next() {
		case 'r':
			h.start = h.pos - 2
			return verbReset
		case 'h':
			h.start = h.pos - 2
			h.pos++ // skip the [
			return scanHighlight
		case eof:
			return nil
		}
	}
}

// verbReset replaces the reset verb with the reset control sequence.
func verbReset(h *highlighter) stateFn {
	h.attrs = []byte(attr["reset"])
	h.replace()
	return scanText
}

// scanHighlight scans the highlight verb for attributes,
// then replaces it with a control sequence derived from said attributes.
func scanHighlight(h *highlighter) stateFn {
	for {
		r := h.next()
		switch {
		case r == eof:
			return nil
		case r == ']':
			if h.attrs != nil {
				h.replace()
			}
			return scanText
		case r == '+':
			continue
		case unicode.IsLetter(r):
			h.backup()
			return scanAttribute
		default:
			h.backup()
			return scanText
		}
	}
}

// scans a attribute and adds it to h.attrs.
func scanAttribute(h *highlighter) stateFn {
	r := h.next()
	switch {
	case r == eof:
		return nil
	case r == 'f' || r == 'b':
		if h.next() == 'g' {
			if unicode.IsNumber(h.next()) {
				h.pos -= 3
				return scanColor256
			}
			h.pos--
		}
		h.backup()
		fallthrough
	case unicode.IsLetter(r):
		start := h.pos - 1
		for {
			r := h.next()
			switch {
			case r == eof:
				return nil
			case unicode.IsLetter(r):
			default:
				if a, ok := attr[string(h.b[start:h.pos-h.width])]; ok {
					h.attrs = append(h.attrs, []byte(a)...)
				}
				h.backup()
				return scanHighlight
			}
		}
	default:
		h.backup()
		return scanHighlight
	}
}

// scans a 256 color attribute and adds it to h.attrs.
func scanColor256(h *highlighter) stateFn {
	var pre byte
	switch h.next() {
	case 'f':
		pre = '3'
	case 'b':
		pre = '4'
	}
	h.next() // skip the g, in "fg" or "bg"
	// can set here because already know it is a number
	h.next()
	start := h.pos - 1
	for {
		r := h.next()
		switch {
		case r == eof:
			return nil
		case unicode.IsNumber(r):
		default:
			b := make([]byte, len(h.attrs)+6+(h.pos-start))
			copy(b, h.attrs)
			l := len(h.attrs)
			copy(b[l:], []byte{';', pre, '8', ';', '5', ';'})
			l += 6
			copy(b[l:], h.b[start:h.pos-1])
			h.attrs = b
			h.backup()
			return scanHighlight
		}
	}
}
