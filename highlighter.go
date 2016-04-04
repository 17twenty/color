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
	b     string // string being scanned
	pos   int    // position in s
	width int    // width of last rune read from s
	start int    // start position of current verb
	attrs string // attributes of current highlight verb
}

// Highlight replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string
func Highlight(s string) string {
	h := &highlighter{b: s}
	h.run()
	return h.b
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
	r, w := utf8.DecodeRuneInString(h.b[h.pos:])
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
	if h.attrs == "" {
		return
	}
	back := h.pos - h.start
	h.attrs = h.attrs[1:]
	h.b = h.b[:h.pos-back] + csi + h.attrs + "m" + h.b[h.pos:]
	h.pos += len(csi) + len(h.attrs) - back
	h.attrs = ""
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
	h.attrs = attr["reset"]
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
			h.replace()
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
				if a, ok := attr[h.b[start:h.pos-h.width]]; ok {
					h.attrs += a
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
	var pre string
	switch h.next() {
	case 'f':
		pre = "3"
	case 'b':
		pre = "4"
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
			h.attrs += ";" + pre + "8;5;" + h.b[start:h.pos-h.width]
			h.backup()
			return scanHighlight
		}
	}
}
