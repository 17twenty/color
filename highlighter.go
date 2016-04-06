package color

import "unicode"

var (
	noVerbBytes  = []byte("%!(NOVERB)")
	missingBytes = []byte("(MISSING)")
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// highlighter holds the state of the scanner.
type highlighter struct {
	buf   []byte // buffer for result
	s     string // string being scanned
	pos   int    // position in buf
	last  int
	attrs []byte // attributes of current highlight verb
}

// Highlight replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string
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
	if h.pos > h.last+1 {
		h.appendLast()
	}
}

func (h *highlighter) get() rune {
	return rune(h.s[h.pos])
}

func (h *highlighter) appendAttrs() {
	b := make([]byte, len(h.buf)+len(csi)+len(h.attrs)) // leading ';', trailing 'm'
	copy(b, h.buf)
	l := len(h.buf)
	copy(b[l:], csi)
	l += len(csi)
	copy(b[l:], h.attrs[1:]) // ignore ;
	b[len(b)-1] = 'm'
	h.buf = b
}

func (h *highlighter) appendLast() {
	h.buf = append(h.buf, h.s[h.last:h.pos]...)
}

// scanText scans until the next highlight or reset verb.
func scanText(h *highlighter) stateFn {
	h.last = h.pos
	for ; h.pos < len(h.s); h.pos++ {
		if h.get() != '%' {
			continue
		}
		if h.pos > h.last {
			h.appendLast()
			h.last = h.pos
		}
		h.pos++
		if h.pos >= len(h.s) {
			h.buf = append(h.buf, noVerbBytes...)
			return nil
		}
		switch h.get() {
		case 'r':
			return verbReset
		case 'h':
			h.last = h.pos - 1
			h.pos += 2
			return scanHighlight
		}
	}
	return nil
}

// verbReset replaces the reset verb with the reset control sequence.
func verbReset(h *highlighter) stateFn {
	h.attrs = attr["reset"]
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
			return scanAttribute
		case r == '+':
			// skip
		case r == ']':
			h.pos++
			if h.attrs != nil {
				h.appendAttrs()
			} else {
				h.appendLast()
			}
			h.attrs = nil
			return scanText
		default:
			h.attrs = nil
			h.appendLast()
			return scanText
		}
	}
	return nil
}

// scanAttribute scans a named attribute
func scanAttribute(h *highlighter) stateFn {
	start := h.pos
	for ; h.pos < len(h.s); h.pos++ {
		r := h.get()
		switch {
		case unicode.IsLetter(r):
			// continue
		default:
			if a, ok := attr[h.s[start:h.pos]]; ok {
				h.attrs = append(h.attrs, a...)
			}
			return scanHighlight
		}
	}
	return nil
}

// scanColor256 scans a 256 color attribute
func scanColor256(h *highlighter, pre []byte) stateFn {
	h.pos++
	if h.get() != 'g' {
		h.pos--
		return scanAttribute
	}
	h.pos++
	if !unicode.IsNumber(h.get()) {
		h.pos -= 2
		return scanAttribute
	}
	start := h.pos
	for ; h.pos < len(h.s); h.pos++ {
		switch {
		case unicode.IsNumber(h.get()):
			// continue
		default:
			a := h.s[start:h.pos]
			b := make([]byte, len(h.attrs)+len(pre)+len(a))
			copy(b, h.attrs)
			l := len(h.attrs)
			copy(b[l:], pre)
			l += len(pre)
			copy(b[l:], a)
			h.attrs = b
			return scanHighlight
		}
	}
	return nil
}
