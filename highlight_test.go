package color

import (
	"fmt"
	"testing"

	"github.com/nhooyr/terminfo/caps"
)

func exp(s string) string {
	if tiErr != nil {
		return ""
	}
	return s
}

func expF(f string, s string) string {
	if tiErr != nil {
		return s
	}
	return fmt.Sprintf(f, s)
}

func TestModes(t *testing.T) {
	t.Parallel()
	for k, v := range modes {
		exp := expF(ti.Strings[v]+"%s"+ti.Strings[caps.ExitAttributeMode], "hi")
		r := Highlight(fmt.Sprintf("%%h[%s]hi%%r", k))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

func TestColors(t *testing.T) {
	t.Parallel()
	for k, v := range colors {
		exp := expF(ti.Color(v, -1)+"%s", "hi")
		r := Highlight(fmt.Sprintf("%%h[fg%s]hi", k))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		exp = expF(ti.Color(-1, v)+"%s", "hi")
		r = Highlight(fmt.Sprintf("%%h[bg%s]hi", k))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

func TestColors256(t *testing.T) {
	t.Parallel()
	for i := 0; i < 256; i++ {
		exp := expF(ti.Color(i, -1)+"%s"+ti.Strings[caps.ExitAttributeMode], "hi")
		r := Highlight(fmt.Sprintf("%%h[fg%d]hi%%r", i))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		exp = expF(ti.Color(-1, i)+"%s"+ti.Strings[caps.ExitAttributeMode], "hi")
		r = Highlight(fmt.Sprintf("%%h[bg%d]hi%%r", i))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

var combinations = map[string]string{
	"%h[fgRed+bgBlue+bold+underline+fg23+bg235]hi":         expF(ti.Color(caps.Red, caps.Blue)+ti.Strings[caps.EnterBoldMode]+ti.Strings[caps.EnterUnderlineMode]+ti.Color(23, 235)+"%s", "hi"),
	"%h[bgBlue+fgYellow+fgGreen+fg34+blink+dim+reverse]hi": expF(ti.Color(-1, caps.Blue)+ti.Color(caps.Yellow, -1)+ti.Color(caps.Green, -1)+ti.Color(34, -1)+ti.Strings[caps.EnterBlinkMode]+ti.Strings[caps.EnterDimMode]+ti.Strings[caps.EnterReverseMode]+"%s", "hi"),
}

func TestCombinations(t *testing.T) {
	t.Parallel()
	for k, v := range combinations {
		if r := Highlight(k); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

var highlightEdgeCases = map[string]string{
	"%h[fgBrightBlack+%h[fgBlue]": exp(ti.Color(caps.BrightBlack, -1)) + errBadAttr,
	"%h[":                   errShort,
	"%h[f":                  errShort,
	"%h[fg":                 errShort,
	"%h{":                   errInvalid,
	"%h[]":                  errMissing,
	"%%h[fgRed]":            "%%h[fgRed]",
	"%[bg232]":              "%[bg232]",
	"%h[fg132":              errShort,
	"%h[fgMagenta[]":        errBadAttr,
	"%h[fgGreen+lold[]":     exp(ti.Color(caps.Green, -1)) + errBadAttr,
	"%h[fgYellow+%#bgBlue]": exp(ti.Color(caps.Yellow, -1)) + errBadAttr,
	"%h][fgRed+%#bgBlue]":   errInvalid,
	"%h[fgRed+":             exp(ti.Color(caps.Red, -1)) + errShort,
	"%%h%h[fgRed]%%":        "%%h" + exp(ti.Color(caps.Red, -1)) + "%%",
	"%h[dsadadssadas]":      errBadAttr,
	"%":                     "%",
	"%h[fgsadas]":           errBadAttr,
	"%h[fgCyan+%h[bgBlue]":  exp(ti.Color(caps.Cyan, -1)) + errBadAttr,
	"lmaokai":               "lmaokai",
	"%h[fgRed]%h[]":         exp(ti.Color(caps.Red, -1)) + errMissing,
	"%h[bgGjo]%h[bgGreen]":  errBadAttr,
	"%h[fg23a]":             errBadAttr,
}

func TestHighlightEdgeCases(t *testing.T) {
	t.Parallel()
	for k, v := range highlightEdgeCases {
		if r := Highlight(k); r != v {
			t.Errorf("Expected %q from %q but result was %q", v, k, r)
		}
	}
}

var stripEdgeCases = map[string]string{
	"%h[fgRed]%smao%r": "%smao",
	"%":                "%",
	"%c":               "%c",
}

func TestStripEdgeCases(t *testing.T) {
	t.Parallel()
	for k, v := range stripEdgeCases {
		if r := Strip(k); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

const s = `%h[fgBlack]hi%r
%h[fgRed]hi%r
%h[bgGreen]hi%r
%h[bgYellow]hi%r
%h[fgBlue]hi%r
%h[fgMagenta]hi%r
%h[bgCyan]hi%r
%h[bgWhite]hi%r
%h[bold]hi%r
%h[underline]hi%r
%h[blink]hi%r
%h[fg22]hi%r
%h[fg233]hi%r
%h[bg3]hi%r
%h[bg102]hi%r
%h[fgBrightBlack+bgCyan+bold+underline+fg32+bg69]hi%r
%h[fg32+bg123+bold+underline+bgBlue+fgBrightGreen+bgWhite]hi%r`

func BenchmarkHighlight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Highlight(s)
	}
}

func BenchmarkStrip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Strip(s)
	}
}
