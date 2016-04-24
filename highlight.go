package color

import (
	"bytes"
	"io"
	"strconv"
	"sync"
	"unicode"

	"github.com/nhooyr/terminfo"
	"github.com/nhooyr/terminfo/caps"
)

const (
	errInvalid = "%%!h(INVALID)" // invalid character in the verb
	errMissing = "%%!h(MISSING)" // no attributes in the verb
	errShort   = "%%!h(SHORT)"   // string ended before the verb
	errBadAttr = "%%!h(BADATTR)" // unknown attribute in the verb
)

// highlighter holds the state of the scanner.
type highlighter struct {
	s     string        // string being scanned
	pos   int           // position in s
	buf   *bytes.Buffer // where result is built
	color bool          // color or strip the highlight verbs
	fg    bool          // foreground or background color attribute
}

// Global terminfo struct.
// TODO no global pls.
var ti, tiErr = terminfo.OpenEnv()

// Highlight replaces the highlight verbs in s with the appropriate control sequences and
// then returns the resulting string.
// It is a thin wrapper around Run.
func Highlight(s string) string {
	return Run(s, true)
}

// Strip removes all highlight verbs in s and then returns the resulting string.
// It is a thin wrapper around Run.
func Strip(s string) string {
	return Run(s, false)
}

// Run runs a highlighter with s as the input and then returns the output. The color argument
// determines whether the highlight verbs will be replaced with their appropriate control
// sequences or instead stripped.
func Run(s string, color bool) string {
	hl := getHighlighter(s, color)
	defer hl.free()
	return hl.run()
}

// highlighterPool allows the reuse of highlighters to avoid allocations.
var highlighterPool = sync.Pool{
	New: func() interface{} {
		hl := new(highlighter)
		// The initial capacity avoids constant reallocation during growth.
		hl.buf = bytes.NewBuffer(make([]byte, 0, 45))
		return hl
	},
}

// getHighlighter returns a new initialized highlighter from the pool.
func getHighlighter(s string, color bool) (hl *highlighter) {
	hl = highlighterPool.Get().(*highlighter)
	hl.s = s
	if tiErr != nil {
		color = false
	}
	hl.color = color
	return
}

// free resets the highlighter.
func (hl *highlighter) free() {
	hl.buf.Reset()
	hl.pos = 0
	highlighterPool.Put(hl)
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*highlighter) stateFn

// run runs the state machine for the highlighter.
func (hl *highlighter) run() string {
	for state := scanText; state != nil; {
		state = state(hl)
	}
	return hl.buf.String()
}

// get returns the current character.
func (hl *highlighter) get() (byte, error) {
	if hl.pos >= len(hl.s) {
		return 0, io.EOF
	}
	return hl.s[hl.pos], nil
}

// writePrev writes n previous bytes to the buffer.
func (hl *highlighter) writePrev(n int) {
	hl.buf.WriteString(hl.s[hl.pos-n : hl.pos])
}

// writeFrom writes the bytes from ppos to pos to the buffer.
func (hl *highlighter) writeFrom(ppos int) {
	if hl.pos > ppos {
		hl.buf.WriteString(hl.s[ppos:hl.pos])
	}
}

func (hl *highlighter) writeAttr(a string) {
	if hl.color {
		hl.buf.WriteString(a)
	}
}

// scanAttribute returns the string from the current character to
// the start of the next attribute or end of the verb.
func (hl *highlighter) scanAttribute() (string, error) {
	start := hl.pos
	hl.pos++
	for {
		ch, err := hl.get()
		if err != nil {
			return "", err
		}
		if ch == '+' || ch == ']' {
			break
		}
		hl.pos++
	}
	return hl.s[start:hl.pos], nil
}

// scanText scans until the next verb.
func scanText(hl *highlighter) stateFn {
	ppos := hl.pos
	// Find next verb.
	for {
		ch, err := hl.get()
		if err != nil {
			hl.writeFrom(ppos)
			return nil
		}
		if ch == '%' {
			hl.writeFrom(ppos)
			hl.pos++
			return scanVerb
		}
		hl.pos++
	}
}

// scanVerb scans the current verb.
func scanVerb(hl *highlighter) stateFn {
	ch, err := hl.get()
	if err != nil {
		// Let fmt insert "%!h(NOVERB)".
		hl.buf.WriteByte('%')
		return nil
	}
	hl.pos++
	switch ch {
	case 'r':
		hl.writeAttr(ti.StringCaps[caps.ExitAttributeMode])
		return scanText
	case 'h':
		// Ensure next character is '['.
		ch, err = hl.get()
		if err != nil {
			hl.buf.WriteString(errShort)
			return nil
		}
		if ch != '[' {
			hl.buf.WriteString(errInvalid)
			return nil
		}
		// Ensure next character is not ']'.
		hl.pos++
		ch, err = hl.get()
		if err != nil {
			hl.buf.WriteString(errShort)
			return nil
		}
		if ch == ']' {
			hl.buf.WriteString(errMissing)
			return nil
		}
		return startAttribute
	}
	// Include the verb.
	hl.writePrev(2)
	return scanText
}

// startAttribute checks the type of the attribute and passes control appropriately.
func startAttribute(hl *highlighter) stateFn {
	// No need to check error because the character was already read.
	switch ch, _ := hl.get(); ch {
	case 'f':
		hl.fg = true
	case 'b':
		hl.fg = false
	default:
		return scanMode
	}
	// Attribute starts with 'f' or 'b' so it could be a color attribute.
	// Rest of the code confirms if it is a color attribute, and if so,
	// whether it is a named or a 256 color attribute.
	hl.pos++
	ch, err := hl.get()
	if err != nil {
		hl.buf.WriteString(errShort)
		return nil
	}
	if ch != 'g' {
		// Actually a mode attribute, because color attributes must begin with "fg" or "bg".
		hl.pos--
		return scanMode
	}
	// Now check if it is a named or 256 color attribute.
	hl.pos++
	ch, err = hl.get()
	if err != nil {
		hl.buf.WriteString(errShort)
		return nil
	}
	if unicode.IsNumber(rune(ch)) {
		return scanColor256
	}
	return scanColor
}

// modes maps mode names to their string capacity positions.
var modes = map[string]int{
	"reset":     caps.ExitAttributeMode,
	"bold":      caps.EnterBoldMode,
	"underline": caps.EnterUnderlineMode,
	"reverse":   caps.EnterReverseMode,
	"blink":     caps.EnterBlinkMode,
	"dim":       caps.EnterDimMode,
}

// scanAttribute scans a mode attribute.
func scanMode(hl *highlighter) stateFn {
	a, err := hl.scanAttribute()
	if err != nil {
		hl.buf.WriteString(errShort)
		return nil
	}
	if n, ok := modes[a]; ok {
		hl.writeAttr(ti.StringCaps[n])
		return endAttribute
	}
	hl.buf.WriteString(errBadAttr)
	return nil
}

// colors maps color names to their integer values.
var colors = map[string]int{
	"Black":         caps.Black,
	"Red":           caps.Red,
	"Green":         caps.Green,
	"Yellow":        caps.Yellow,
	"Blue":          caps.Blue,
	"Magenta":       caps.Magenta,
	"Cyan":          caps.Cyan,
	"White":         caps.White,
	"BrightBlack":   caps.BrightBlack,
	"BrightRed":     caps.BrightRed,
	"BrightGreen":   caps.BrightGreen,
	"BrightYellow":  caps.BrightYellow,
	"BrightBlue":    caps.BrightBlue,
	"BrightMagenta": caps.BrightMagenta,
	"BrightCyan":    caps.BrightCyan,
	"BrightWhite":   caps.BrightWhite,
}

// scanColor scans a named color attribute.
func scanColor(hl *highlighter) stateFn {
	a, err := hl.scanAttribute()
	if err != nil {
		hl.buf.WriteString(errShort)
		return nil
	}
	if c, ok := colors[a]; ok {
		if hl.fg {
			hl.writeAttr(ti.Color(c, -1))
		} else {
			hl.writeAttr(ti.Color(-1, c))
		}
		return endAttribute
	}
	hl.buf.WriteString(errBadAttr)
	return nil
}

// scanColor256 scans a 256 color attribute.
func scanColor256(hl *highlighter) stateFn {
	a, err := hl.scanAttribute()
	if err != nil {
		hl.buf.WriteString(errShort)
		return nil
	}
	t, err := strconv.Atoi(a)
	if err != nil {
		hl.buf.WriteString(errBadAttr)
		return nil
	}
	if hl.fg {
		hl.writeAttr(ti.Color(t, -1))
	} else {
		hl.writeAttr(ti.Color(-1, t))
	}
	return endAttribute
}

// endAttribute handles the end of attributes. If there is another attribute, control is
// thrown to scanHighlight, but if the verb has ended, control is thrown to scanText.
func endAttribute(hl *highlighter) stateFn {
	ch, _ := hl.get()
	hl.pos++
	if ch == ']' {
		return scanText
	}
	// Must read the next character here because scanHighlight assumes that
	// the character was already read. See scanVerb.
	if _, err := hl.get(); err != nil {
		hl.buf.WriteString(errShort)
		return nil
	}
	return startAttribute
}
