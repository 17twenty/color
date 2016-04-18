package color

import (
	"strings"
	"sync"
	"unicode"
)

const (
	errMissing = "%%!h(MISSING)" // no attributes in the highlight verb
	errInvalid = "%%!h(INVALID)" // invalid character in the highlight verb
	errBadAttr = "%%!h(BADATTR)" // unknown attribute in the highlight verb
)

// highlighter holds the state of the scanner.
type highlighter struct {
	s     string // string being scanned
	pos   int    // position in s
	buf   buffer // where result is built
	attrs buffer // attributes of current verb
	color bool   // color or strip the highlight verbs
}

// Highlight replaces the highlight verbs in s with the appropriate control sequences and
// then returns the resulting string.
// It is a thin wrapper around color.Run().
func Highlight(s string) string {
	return Run(s, true)
}

// Strip removes all highlight verbs in s and then returns the resulting string.
// It is a thin wrapper around color.Run().
func Strip(s string) string {
	return Run(s, false)
}

// Run runs a highlighter with s as the input and then returns the output. The strip argument
// determines whether the highlight verbs will be stripped or instead replaced with
// their appropriate control sequences.
// Do not use this directly unless you know what you are doing.
func Run(s string, color bool) string {
	hl := getHighlighter(s, color)
	defer hl.free()
	hl.run()
	s = string(hl.buf)
	return s
}

// highlighterPool allows the reuse of highlighters to avoid allocations.
var highlighterPool = sync.Pool{
	New: func() interface{} {
		hl := new(highlighter)
		// The initial capacities avoid constant reallocation during growth.
		hl.buf = make([]byte, 0, 45)
		hl.attrs = make([]byte, 0, 15)
		return hl
	},
}

// getHighlighter returns a new initialized highlighter from the pool.
func getHighlighter(s string, color bool) (hl *highlighter) {
	hl = highlighterPool.Get().(*highlighter)
	hl.s, hl.color = s, color
	return
}

// free resets the highlighter.
func (hl *highlighter) free() {
	hl.buf.reset()
	hl.pos = 0
	highlighterPool.Put(hl)
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// run runs the state machine for the highlighter.
func (hl *highlighter) run() {
	for state := scanText; state != nil; {
		state = state(hl)
	}
}

const eof = -1

// get returns the current rune.
func (hl *highlighter) get() rune {
	if hl.pos >= len(hl.s) {
		return eof
	}
	return rune(hl.s[hl.pos])
}

// writeAttrs writes a control sequence derived from h.attrs[1:] to h.buf.
func (hl *highlighter) writeAttrs() {
	hl.buf.writeString(csi)
	hl.buf.write(hl.attrs[1:])
	hl.buf.writeByte('m')
}

// writePrev writes n previous characters to the buffer
func (hl *highlighter) writePrev(n int) {
	hl.pos++
	hl.buf.writeString(hl.s[hl.pos-n : hl.pos])
}

// writeFrom writes the characters from ppos to pos to the buffer.
func (hl *highlighter) writeFrom(ppos int) {
	if hl.pos > ppos {
		// Append remaining characters.
		hl.buf.writeString(hl.s[ppos:hl.pos])
	}
}

// scanText scans until the next verb.
func scanText(hl *highlighter) stateFn {
	ppos := hl.pos
	// Find next verb.
	for {
		switch hl.get() {
		case eof:
			hl.writeFrom(ppos)
			return nil
		case '%':
			hl.writeFrom(ppos)
			hl.pos++
			if hl.color {
				return scanVerb
			}
			return stripVerb
		}
		hl.pos++
	}
}

// stripVerb skips the current verb.
func stripVerb(hl *highlighter) stateFn {
	switch hl.get() {
	case 'r':
		// Strip the reset verb.
		hl.pos++
	case 'h':
		// Strip inside the highlight verb.
		hl.pos++
		j := strings.IndexByte(hl.s[hl.pos:], ']')
		if j == -1 {
			hl.buf.writeString(errInvalid)
			return nil
		}
		hl.pos += j + 1
	case eof:
		// Let fmt handle "%!h(NOVERB)".
		hl.buf.writeByte('%')
		return nil
	default:
		// Include the verb.
		hl.writePrev(2)
	}
	return scanText
}

// scanVerb scans the current verb.
func scanVerb(hl *highlighter) stateFn {
	switch hl.get() {
	case 'r':
		hl.pos++
		return verbReset
	case 'h':
		hl.pos += 2
		return scanHighlight
	case eof:
		// Let fmt handle "%!h(NOVERB)".
		hl.buf.writeByte('%')
		return nil
	}
	hl.writePrev(2)
	return scanText
}

// verbReset writes the reset verb with the reset control sequence.
func verbReset(hl *highlighter) stateFn {
	hl.attrs.writeString(attrs["reset"])
	hl.writeAttrs()
	hl.attrs.reset()
	return scanText
}

// scanHighlight scans the highlight verb for attributes,
// then writes a control sequence derived from said attributes to the buffer.
func scanHighlight(hl *highlighter) stateFn {
	r := hl.get()
	switch {
	case r == 'f':
		return scanColor256(hl, preFg256)
	case r == 'b':
		return scanColor256(hl, preBg256)
	case unicode.IsLetter(r):
		return scanAttribute(hl, 0)
	case r == '+':
		hl.pos++
		return scanHighlight
	case r == ']':
		if len(hl.attrs) != 0 {
			hl.writeAttrs()
		} else {
			hl.buf.writeString(errMissing)
		}
		hl.attrs.reset()
		hl.pos++
		return scanText
	default:
		return abortHighlight(hl, errInvalid)
	}
}

// scanAttribute scans a named attribute.
func scanAttribute(hl *highlighter, off int) stateFn {
	start := hl.pos - off
	for unicode.IsLetter(hl.get()) {
		hl.pos++
	}
	if a, ok := attrs[hl.s[start:hl.pos]]; ok {
		hl.attrs.writeString(a)
	} else {
		return abortHighlight(hl, errBadAttr)
	}
	return scanHighlight
}

// abortHighlight writes a error to the buffer and
// then skips to the end of the highlight verb.
func abortHighlight(hl *highlighter, msg string) stateFn {
	hl.buf.writeString(msg)
	hl.attrs.reset()
	for {
		switch hl.get() {
		case ']':
			hl.pos++
			return scanText
		case eof:
			return nil
		}
		hl.pos++
	}
}

// scanColor256 scans a 256 color attribute.
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
	hl.pos++
	for unicode.IsNumber(hl.get()) {
		hl.pos++
	}
	hl.attrs.writeString(pre)
	hl.attrs.writeString(hl.s[start:hl.pos])
	return scanHighlight
}
