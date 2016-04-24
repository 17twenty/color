package color

import (
	"bytes"
	"io"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/nhooyr/terminfo"
	"github.com/nhooyr/terminfo/caps"
)

const (
	errInvalid = "%%!h(INVALID)" // invalid character in the highlight verb
	errBadAttr = "%%!h(BADATTR)" // unknown attribute in the highlight verb
)

var colors = map[string]int{
	"Black":   terminfo.ColorBlack,
	"Maroon":  terminfo.ColorMaroon,
	"Green":   terminfo.ColorGreen,
	"Olive":   terminfo.ColorOlive,
	"Navy":    terminfo.ColorNavy,
	"Purple":  terminfo.ColorPurple,
	"Teal":    terminfo.ColorTeal,
	"Silver":  terminfo.ColorSilver,
	"Gray":    terminfo.ColorGray,
	"Red":     terminfo.ColorRed,
	"Lime":    terminfo.ColorLime,
	"Yellow":  terminfo.ColorYellow,
	"Blue":    terminfo.ColorBlue,
	"Fuchsia": terminfo.ColorFuchsia,
	"Aqua":    terminfo.ColorAqua,
	"White":   terminfo.ColorWhite,
}

// highlighter holds the state of the scanner.
type highlighter struct {
	s       string        // string being scanned
	pos     int           // position in s
	buf     *bytes.Buffer // where result is built
	color   bool          // color or strip the highlight verbs
	fg      bool          // foreground or background color attribute
	noAttrs bool          // not written attrs to buf
}

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

// Run runs a highlighter with s as the input and then returns the output. The strip argument
// determines whether the highlight verbs will be stripped or instead replaced with
// their appropriate control sequences.
// Do not use this directly unless you know what you are doing.
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
	if tiErr == nil {
		hl.color = color
	} else {
		hl.color = false
	}
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

// get returns the current byte.
func (hl *highlighter) get() (byte, error) {
	if hl.pos >= len(hl.s) {
		return 0, io.EOF
	}
	return hl.s[hl.pos], nil
}

// writePrev writes n previous characters to the buffer
func (hl *highlighter) writePrev(n int) {
	hl.pos++
	hl.buf.WriteString(hl.s[hl.pos-n : hl.pos])
}

// writeFrom writes the characters from ppos to pos to the buffer.
func (hl *highlighter) writeFrom(ppos int) {
	if hl.pos > ppos {
		// Append remaining characters.
		hl.buf.WriteString(hl.s[ppos:hl.pos])
	}
}

func (hl *highlighter) scanAttribute() (string, error) {
	start := hl.pos
	hl.pos++
	for {
		ch, err := hl.get()
		if err != nil {
			hl.buf.WriteString(errInvalid)
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
		} else if ch == '%' {
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
	ch, err := hl.get()
	if err != nil {
		// Let fmt handle "%!h(NOVERB)".
		hl.buf.WriteByte('%')
		return nil
	}
	switch ch {
	case 'r':
		// Strip the reset verb.
		hl.pos++
	case 'h':
		// Strip inside the highlight verb.
		hl.pos++
		j := strings.IndexByte(hl.s[hl.pos:], ']')
		if j == -1 {
			hl.buf.WriteString(errInvalid)
			return nil
		}
		hl.pos += j + 1
	default:
		// Include the verb.
		hl.writePrev(2)
	}
	return scanText
}

// scanVerb scans the current verb.
func scanVerb(hl *highlighter) stateFn {
	ch, err := hl.get()
	if err != nil {
		// Let fmt handle "%!h(NOVERB)".
		hl.buf.WriteByte('%')
		return nil
	}
	switch ch {
	case 'r':
		hl.pos++
		return verbReset
	case 'h':
		hl.pos++
		ch, err = hl.get()
		if err != nil || ch != '[' {
			hl.buf.WriteString(errInvalid)
			return nil
		}
		hl.noAttrs = true
		hl.pos++
		return scanHighlight
	}
	hl.writePrev(2)
	return scanText
}

// verbReset writes the reset verb with the reset control sequence.
func verbReset(hl *highlighter) stateFn {
	hl.buf.WriteString(ti.StringCaps[caps.ExitAttributeMode])
	return scanText
}

// scanHighlight scans the highlight verb for attributes,
// then writes a control sequence derived from said attributes to the buffer.
func scanHighlight(hl *highlighter) stateFn {
	ch, err := hl.get()
	if err != nil {
		hl.buf.WriteString(errInvalid)
		return nil
	}
	switch {
	case ch == 'f':
		hl.fg = true
		return scanColor
	case ch == 'b':
		hl.fg = false
		return scanColor
	default:
		return scanMode
	}
}

// scanAttribute scans a mode.
func scanMode(hl *highlighter) stateFn {
	a, err := hl.scanAttribute()
	if err != nil {
		return nil
	}
	switch a {
	case "bold":
		a = ti.StringCaps[caps.EnterBoldMode]
	case "underline":
		a = ti.StringCaps[caps.EnterUnderlineMode]
	case "reverse":
		a = ti.StringCaps[caps.EnterReverseMode]
	case "blink":
		a = ti.StringCaps[caps.EnterBlinkMode]
	case "dim":
		a = ti.StringCaps[caps.EnterDimMode]
	case "reset":
		a = ti.StringCaps[caps.ExitAttributeMode]
	default:
		hl.buf.WriteString(errBadAttr)
		return nil
	}
	hl.buf.WriteString(a)
	return successAttribute
}

// scanColor scans a color attribute.
func scanColor(hl *highlighter) stateFn {
	hl.pos++
	ch, err := hl.get()
	if err != nil {
		return nil
	}
	if ch != 'g' {
		hl.pos--
		return scanMode
	}
	hl.pos++
	ch, err = hl.get()
	if err != nil {
		return nil
	}
	if unicode.IsNumber(rune(ch)) {
		return scanColor256
	}
	a, err := hl.scanAttribute()
	if err != nil {
		return nil
	}
	if c, ok := colors[a]; ok {
		if hl.fg {
			hl.buf.WriteString(ti.Color(c, -1))
		} else {
			hl.buf.WriteString(ti.Color(-1, c))
		}
		return successAttribute
	}
	hl.buf.WriteString(errBadAttr)
	return nil
}

// scanColor256 scans a 256 color attribute.
func scanColor256(hl *highlighter) stateFn {
	a, err := hl.scanAttribute()
	if err != nil {
		return nil
	}
	t, _ := strconv.Atoi(a)
	if hl.fg {
		hl.buf.WriteString(ti.Color(t, -1))
	} else {
		hl.buf.WriteString(ti.Color(-1, t))
	}
	return successAttribute
}

func successAttribute(hl *highlighter) stateFn {
	hl.noAttrs = false
	ch, _ := hl.get()
	hl.pos++
	if ch == ']' {
		return scanText
	}
	return scanHighlight
}
