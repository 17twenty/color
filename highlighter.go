package color

import (
	"unicode"
	"unicode/utf8"
)

const eof = -1

type stateFn func(*highlighter) stateFn

type highlighter struct {
	s              string
	pos            int
	width          int
	startHighlight int
	codes          string
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

func (h *highlighter) replace(back int) {
	h.s = h.s[:h.pos-back] + csi + h.codes[:len(h.codes)-1] + "m" + h.s[h.pos:]
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
		case 'h':
			h.pos -= 2 // backup to %
			h.startHighlight = h.pos
			h.pos += 3 // skip the %h#
			return scanHighlight
		case 'r':
			return verbReset
		case eof:
			return nil
		}
	}
}

// replaces the highlight verb with the appropiate control sequence
func scanHighlight(h *highlighter) stateFn {
	for {
		r := h.next()
		switch {
		case r == eof:
			return nil
		case r == '#':
			h.replace(h.pos - h.startHighlight)
		case r == '+':
			continue
		case unicode.IsLetter(r):
			h.backup()
			return scanAttribute
		default:
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
				h.codes += a
			}
			h.backup()
			return scanHighlight
		}
	}
}

func scanColor256(h *highlighter) stateFn {
	var b, prefix string
	switch string(h.next()) + string(h.next()) {
	case "fg":
		prefix = "3"
	case "bg":
		prefix = "4"
	}
	for {
		r := h.next()
		switch {
		case r == eof:
			return nil
		case unicode.IsNumber(r):
			b += string(r)
		default:
			if b != "" {
				h.codes += prefix + "8;5;" + b + ";"
			}
			h.backup()
			return scanHighlight
		}
	}
}

// replaces the reset verb with the reset control sequence
func verbReset(h *highlighter) stateFn {
	h.codes = attrs["reset"]
	h.replace(2)
	return scanText
}
