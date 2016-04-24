package color_test

import (
	"fmt"
	"testing"

	"github.com/nhooyr/color"
	"github.com/nhooyr/terminfo/caps"
)

var ti, tiErr = color.Ti, color.TiErr

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
	for k, v := range color.Modes {
		exp := expF(ti.StringCaps[v]+"%s"+ti.StringCaps[caps.ExitAttributeMode], "hi")
		r := color.Highlight(fmt.Sprintf("%%h[%s]hi%%r", k))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

func TestColors(t *testing.T) {
	for k, v := range color.Colors {
		exp := expF(ti.Color(v, -1)+"%s", "hi")
		r := color.Highlight(fmt.Sprintf("%%h[fg%s]hi", k))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		exp = expF(ti.Color(-1, v)+"%s", "hi")
		r = color.Highlight(fmt.Sprintf("%%h[bg%s]hi", k))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

func TestColors256(t *testing.T) {
	for i := 0; i < 256; i++ {
		exp := expF(ti.Color(i, -1)+"%s"+ti.StringCaps[caps.ExitAttributeMode], "hi")
		r := color.Highlight(fmt.Sprintf("%%h[fg%d]hi%%r", i))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		exp = expF(ti.Color(-1, i)+"%s"+ti.StringCaps[caps.ExitAttributeMode], "hi")
		r = color.Highlight(fmt.Sprintf("%%h[bg%d]hi%%r", i))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

var combinations = map[string]string{
	"%h[fgMaroon+bgNavy+bold+underline+fg23+bg235]hi":     expF(ti.Color(caps.Maroon, caps.Navy)+ti.StringCaps[caps.EnterBoldMode]+ti.StringCaps[caps.EnterUnderlineMode]+ti.Color(23, 235)+"%s", "hi"),
	"%h[bgBlue+fgOlive+fgGreen+fg34+blink+dim+reverse]hi": expF(ti.Color(-1, caps.Blue)+ti.Color(caps.Olive, -1)+ti.Color(caps.Green, -1)+ti.Color(34, -1)+ti.StringCaps[caps.EnterBlinkMode]+ti.StringCaps[caps.EnterDimMode]+ti.StringCaps[caps.EnterReverseMode]+"%s", "hi"),
}

func TestCombinations(t *testing.T) {
	for k, v := range combinations {
		if r := color.Highlight(k); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

var highlightEdgeCases = map[string]string{
	"%h[fgGray+%h[fgBlue]": exp(ti.Color(caps.Gray, -1)) + color.ErrBadAttr,
	"%h[":                  color.ErrShort,
	"%h{":                  color.ErrInvalid,
	"%h[]":                 color.ErrMissing,
	"%%h[fgRed]":           "%%h[fgRed]",
	"%[bg232]":             "%[bg232]",
	"%h[fg132":             color.ErrShort,
	"%h[fgFuchsia[]":       color.ErrBadAttr,
	"%h[fgGreen+lold[]":    exp(ti.Color(caps.Green, -1)) + color.ErrBadAttr,
	"%h[fgOlive+%#bgBlue]": exp(ti.Color(caps.Olive, -1)) + color.ErrBadAttr,
	"%h][fgRed+%#bgBlue]":  color.ErrInvalid,
	"%h[fgRed+":            exp(ti.Color(caps.Red, -1)) + color.ErrShort,
	"%%h%h[fgRed]%%":       "%%h" + exp(ti.Color(caps.Red, -1)) + "%%",
	"%h[dsadadssadas]":     color.ErrBadAttr,
	"%":                    "%",
	"%h[fgsadas]":          color.ErrBadAttr,
	"%h[fgAqua+%h[bgBlue]": exp(ti.Color(caps.Aqua, -1)) + color.ErrBadAttr,
	"lmaokai":              "lmaokai",
	"%h[fgMaroon]%h[]":     exp(ti.Color(caps.Maroon, -1)) + color.ErrMissing,
	"%h[bgGjo]%h[bgGreen]": color.ErrBadAttr,
	"%h[fg23a]":            color.ErrBadAttr,
}

func TestHighlightEdgeCases(t *testing.T) {
	for k, v := range highlightEdgeCases {
		if r := color.Highlight(k); r != v {
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
	for k, v := range stripEdgeCases {
		if r := color.Strip(k); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

const s = `%h[fgBlack]hi%r
%h[fgMaroon]hi%r
%h[bgGreen]hi%r
%h[bgOlive]hi%r
%h[fgBlue]hi%r
%h[fgFuchsia]hi%r
%h[bgAqua]hi%r
%h[bgWhite]hi%r
%h[bold]hi%r
%h[underline]hi%r
%h[blink]hi%r
%h[fg22]hi%r
%h[fg233]hi%r
%h[bg3]hi%r
%h[bg102]hi%r
%h[fgGray+bgAqua+bold+underline+fg32+bg69]hi%r
%h[fg32+bg123+bold+underline+bgNavy+fgLime+bgWhite]hi%r`

var result interface{}

func BenchmarkHighlight(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = color.Highlight(s)
	}
	result = r
}

func BenchmarkStrip(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = color.Strip(s)
	}
	result = r
}
