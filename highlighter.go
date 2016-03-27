package color

import (
	"unicode"
	"unicode/utf8"
)

const eof = -1

type stateFn func(*highlighter) stateFn

type highlighter struct {
	s     string // string to search and replace the verbs
	pos   int    // position in s
	width int    // width of last rune read from s
	start int    // start position of current highlight verb
	codes string // codes of current highlight verb
}

func (h *highlighter) run() {
	for state := scanText; state != nil; {
		state = state(h)
	}
}

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

func (h *highlighter) backup() {
	h.pos -= h.width
}

// replaces previous back characters with h.codes[1:] (remove first semicolon)
func (h *highlighter) replace(back int) {
	h.s = h.s[:h.pos-back] + csi + h.codes[1:] + "m" + h.s[h.pos:]
	h.codes = ""
	h.pos -= back
}

// scans until the next highlight or reset verb
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
			return verbReset
		case 'h':
			h.pos -= 2 // backup to %
			h.start = h.pos
			h.pos += 3 // skip the %h[
			return scanHighlight
		case eof:
			return nil
		}
	}
}

// replaces the reset verb with the reset control sequence
func verbReset(h *highlighter) stateFn {
	h.codes = ";" + attrs["reset"]
	h.replace(2)
	return scanText
}

// replaces the highlight verb with the appropiate control sequence
func scanHighlight(h *highlighter) stateFn {
	for {
		r := h.next()
		switch {
		case r == eof:
			return nil
		case r == ']':
			if h.codes != "" {
				h.replace(h.pos - h.start)
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
			if a, ok := attrs[b]; ok {
				h.codes += ";" + a
			}
			h.backup()
			return scanHighlight
		}
	}
}

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
			if b != "" {
				h.codes += ";" + pre + "8;5;" + b
			}
			h.backup()
			return scanHighlight
		}
	}
}
