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
	s     string // string being scanned
	pos   int    // position in s
	width int    // width of last rune read from s
	start int    // start position of current verb
	attrs string // attributes of current highlight verb
}

// run runs the state machine for the highlighter.
func (h *highlighter) run() {
	for state := scanText; state != nil; {
		state = state(h)
	}
}

// next returns the next rune in the input.
func (h *highlighter) next() rune {
	if h.pos >= len(h.s) {
		h.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(h.s[h.pos:])
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
	h.s = h.s[:h.pos-back] + csi + h.attrs[1:] + "m" + h.s[h.pos:]
	h.attrs = ""
	h.pos -= back
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
	h.attrs = ";" + attr["reset"]
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
			if h.attrs != "" {
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
	var b string
	for {
		r := h.next()
		switch {
		case r == eof:
			return nil
		case unicode.IsLetter(r):
			b += string(r)
		case unicode.IsNumber(r) && (b == "fg" || b == "bg"):
			h.pos -= 3
			return scanColor256
		default:
			if a, ok := attr[b]; ok {
				h.attrs += ";" + a
			}
			h.backup()
			return scanHighlight
		}
	}
}

// scans a 256 color attribute and adds it to h.attrs.
func scanColor256(h *highlighter) stateFn {
	var b, pre string
	switch h.next() {
	case 'f':
		pre = "3"
	case 'b':
		pre = "4"
	}
	h.next() // skip the g, in "fg" or "bg"
	for {
		r := h.next()
		switch {
		case r == eof:
			return nil
		case unicode.IsNumber(r):
			b += string(r)
		default:
			h.attrs += ";" + pre + "8;5;" + b
			h.backup()
			return scanHighlight
		}
	}
}
