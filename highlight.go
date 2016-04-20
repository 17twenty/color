package color

import (
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/gdamore/tcell"
)

const (
	errMissing = "%%!h(MISSING)" // no attributes in the highlight verb
	errInvalid = "%%!h(INVALID)" // invalid character in the highlight verb
	errBadAttr = "%%!h(BADATTR)" // unknown attribute in the highlight verb
)

var colors = map[string]tcell.Color{
	"Black":   tcell.ColorBlack,
	"Maroon":  tcell.ColorMaroon,
	"Green":   tcell.ColorGreen,
	"Olive":   tcell.ColorOlive,
	"Navy":    tcell.ColorNavy,
	"Purple":  tcell.ColorPurple,
	"Teal":    tcell.ColorTeal,
	"Silver":  tcell.ColorSilver,
	"Gray":    tcell.ColorGray,
	"Red":     tcell.ColorRed,
	"Lime":    tcell.ColorLime,
	"Yellow":  tcell.ColorYellow,
	"Blue":    tcell.ColorBlue,
	"Fuchsia": tcell.ColorFuchsia,
	"Aqua":    tcell.ColorAqua,
	"White":   tcell.ColorWhite,
}

// highlighter holds the state of the scanner.
type highlighter struct {
	s     string // string being scanned
	pos   int    // position in s
	buf   buffer // where result is built
	color bool   // color or strip the highlight verbs
	fg    bool   // foreground or background color attribute
	noAttrs bool   // not written attrs to buf
	ti    *tcell.Terminfo
}

// Highlight replaces the highlight verbs in s with the appropriate control sequences and
// then returns the resulting string.
// It is a thin wrapper around Run.
func Highlight(s string, ti *tcell.Terminfo) string {
	return Run(s, true, ti)
}

// Strip removes all highlight verbs in s and then returns the resulting string.
// It is a thin wrapper around Run.
func Strip(s string, ti *tcell.Terminfo) string {
	return Run(s, false, ti)
}

// Run runs a highlighter with s as the input and then returns the output. The strip argument
// determines whether the highlight verbs will be stripped or instead replaced with
// their appropriate control sequences.
// Do not use this directly unless you know what you are doing.
func Run(s string, color bool, ti *tcell.Terminfo) string {
	hl := getHighlighter(s, color, ti)
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
		return hl
	},
}

// getHighlighter returns a new initialized highlighter from the pool.
func getHighlighter(s string, color bool, ti *tcell.Terminfo) (hl *highlighter) {
	hl = highlighterPool.Get().(*highlighter)
	hl.s, hl.color = s, color
	if hl.color {
		hl.ti = ti
	}
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
				hl.noAttrs = true
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
	hl.buf.writeString(hl.ti.AttrOff)
	return scanText
}

// scanHighlight scans the highlight verb for attributes,
// then writes a control sequence derived from said attributes to the buffer.
func scanHighlight(hl *highlighter) stateFn {
	for {
		r := hl.get()
		switch {
		case r == 'f':
			hl.fg = true
			return scanColor
		case r == 'b':
			hl.fg = false
			return scanColor
		case unicode.IsLetter(r):
			return scanAttribute
		case r == '+':
			hl.pos++
			continue
		case r == ']':
			if hl.noAttrs {
				hl.buf.writeString(errMissing)
			}
			hl.pos++
			return scanText
		default:
			return abortHighlight(hl, errInvalid)
		}
	}
}

// abortHighlight writes a error to the buffer and
// then skips to the end of the highlight verb.
func abortHighlight(hl *highlighter, msg string) stateFn {
	hl.buf.writeString(msg)
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

// scanAttribute scans a named attribute.
func scanAttribute(hl *highlighter) stateFn {
	start := hl.pos
	for unicode.IsLetter(hl.get()) {
		hl.pos++
	}
	a := hl.s[start:hl.pos]
	switch a {
	case "bold":
		a = hl.ti.Bold
	case "underline":
		a = hl.ti.Underline
	case "reverse":
		a = hl.ti.Reverse
	case "blink":
		a = hl.ti.Blink
	case "dim":
		a = hl.ti.Dim
	case "attrOff":
		a = hl.ti.AttrOff
	default:
		return abortHighlight(hl, errBadAttr)
	}
	hl.buf.writeString(a)
	hl.noAttrs = false
	return scanHighlight
}

// scanColor scans a color attribute.
func scanColor(hl *highlighter) stateFn {
	hl.pos++
	if hl.get() != 'g' {
		hl.pos--
		return scanAttribute
	}
	hl.pos++
	r := hl.get()
	switch {
	case unicode.IsNumber(r):
		return scanColor256
	case unicode.IsLetter(r):
		// continue
	default:
		return abortHighlight(hl, errBadAttr)
	}
	start := hl.pos
	hl.pos++
	for unicode.IsLetter(hl.get()) {
		hl.pos++
	}
	if c, ok := colors[hl.s[start:hl.pos]]; ok {
		if hl.fg {
			hl.buf.writeString(hl.ti.TColor(c, -1))
		} else {
			hl.buf.writeString(hl.ti.TColor(-1, c))
		}
		hl.noAttrs = false
		return scanHighlight
	}
	return abortHighlight(hl, errBadAttr)
}

// scanColor256 scans a 256 color attribute.
func scanColor256(hl *highlighter) stateFn {
	start := hl.pos
	hl.pos++
	for unicode.IsNumber(hl.get()) {
		hl.pos++
	}
	n, _ := strconv.ParseInt(hl.s[start:hl.pos], 10, 32)
	if hl.fg {
		hl.buf.writeString(hl.ti.TColor(tcell.Color(n), -1))
	} else {
		hl.buf.writeString(hl.ti.TColor(-1, tcell.Color(n)))
	}
	hl.noAttrs = false
	return scanHighlight
}
