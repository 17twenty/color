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
	r     rune   // current rune
	attrs string // attributes of current highlight verb
}

// Highlight replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string
func Highlight(s string) string {
	h := &highlighter{s: s}
	h.run()
	return h.s
}

// run runs the state machine for the highlighter.
func (h *highlighter) run() {
	h.next()
	for state := scanText; state != nil; {
		state = state(h)
	}
}

// next gets the next rune in the string.
func (h *highlighter) next() {
	if h.pos >= len(h.s) {
		h.r = eof
		return
	}
	h.r, h.width = utf8.DecodeRuneInString(h.s[h.pos:])
	h.pos += h.width
}

// replaces the verb with a control sequence derived from h.attrs[1:].
func (h *highlighter) replace() {
	h.attrs = h.attrs[1:]
	h.s = h.s[:h.start] + csi + h.attrs + "m" + h.s[h.pos:]
	h.pos += len(csi) + len(h.attrs) - (h.pos-h.start)
}

// scans until the next highlight or reset verb.
func scanText(h *highlighter) stateFn {
	for ; ; h.next() {
		switch h.r {
		case eof:
			return nil
		case '%':
			// a verb!
		default:
			continue
		}
		h.next()
		switch h.r {
		case 'r':
			h.start = h.pos - 2
			return verbReset
		case 'h':
			h.start = h.pos - 2
			h.pos++ // skip the [
			h.next()
			return scanHighlight
		case eof:
			return nil
		}
	}
}

// verbReset replaces the reset verb with the reset control sequence.
func verbReset(h *highlighter) stateFn {
	h.attrs = attrs["reset"]
	h.replace()
	return scanText
}

// scanHighlight scans the highlight verb for attributes,
// then replaces it with a control sequence derived from said attributes.
func scanHighlight(h *highlighter) stateFn {
	for ; ; h.next() {
		switch {
		case h.r == eof:
			return nil
		case h.r == 'f':
			return scanColor256(h, preFg256)
		case h.r == 'b':
			return scanColor256(h, preBg256)
		case unicode.IsLetter(h.r):
			return scanAttribute
		case h.r == '+':
			// skip
		case h.r == ']':
			if h.attrs != "" {
				h.replace()
			}
			h.next()
			fallthrough
		default:
			h.attrs = ""
			return scanText
		}
	}
}

// scanAttribute scans a named attribute
func scanAttribute(h *highlighter) stateFn {
	start := h.pos - h.width
	for {
		h.next()
		switch {
		case h.r == eof:
			return nil
		case unicode.IsLetter(h.r):
			// continue
		default:
			if a, ok := attrs[h.s[start:h.pos-h.width]]; ok {
				h.attrs += a
			}
			return scanHighlight
		}
	}
}

// scanColor256 scans a 256 color attribute
func scanColor256(h *highlighter, pre string) stateFn {
	h.next()
	if h.r != 'g' {
		h.width++ // start at f/b
		return scanAttribute
	}
	h.next()
	if !unicode.IsNumber(h.r) {
		h.width += 2 // start at (f/b)g
		return scanAttribute
	}
	start := h.pos - h.width
	for {
		h.next()
		switch {
		case h.r == eof:
			return nil
		case unicode.IsNumber(h.r):
			// continue
		default:
			h.attrs += pre + h.s[start:h.pos-h.width]
			return scanHighlight
		}
	}
}
