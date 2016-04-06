package color

import "unicode"

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// highlighter holds the state of the scanner.
type highlighter struct {
	s     string // string being scanned
	pos   int    // position in buf
	start int    // start position of current verb
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
	for state := scanText; state != nil; {
		state = state(h)
	}
}

func (h *highlighter) get() rune {
	return rune(h.s[h.pos])
}

// replaces the verb with a control sequence derived from h.attrs[1:].
func (h *highlighter) replace() {
	h.s = h.s[:h.start] + csi + h.attrs[1:] + "m" + h.s[h.pos:]
	h.pos += len(csi) + len(h.attrs) - (h.pos - h.start)
}

// scanText scans until the next highlight or reset verb.
func scanText(h *highlighter) stateFn {
	for ; h.pos < len(h.s); h.pos++ {
		if h.get() != '%' {
			continue
		}
		h.pos++
		if h.pos >= len(h.s) {
			return nil
		}
		switch h.get() {
		case 'r':
			h.start = h.pos - 1
			h.pos++
			return verbReset
		case 'h':
			h.start = h.pos - 1
			h.pos += 2
			return scanHighlight
		}
	}
	return nil
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
			h.pos++
			if h.attrs != "" {
				h.replace()
			}
			fallthrough
		default:
			h.attrs = ""
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
				h.attrs += a
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
			h.attrs += pre + h.s[start:h.pos]
			return scanHighlight
		}
	}
	return nil
}
