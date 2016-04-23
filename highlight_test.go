package color_test

import (
	"fmt"
	"testing"

	"github.com/nhooyr/color"
	"github.com/nhooyr/terminfo"
)

func TestAttributes(t *testing.T) {
	for k, v := range color.Colors {
		exp := color.Ti.Color(v, -1) + "hi" + color.Ti.Reset
		r := color.Highlight(fmt.Sprintf("%%h[fg%s]hi%%r", k))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		exp = color.Ti.Color(-1, v) + "hi" + color.Ti.Reset
		r = color.Highlight(fmt.Sprintf("%%h[bg%s]hi%%r", k))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

func TestColor256(t *testing.T) {
	for i := 0; i < 256; i++ {
		exp := color.Ti.Color(i, -1) + "hi" + color.Ti.Reset
		r := color.Highlight(fmt.Sprintf("%%h[fg%d]hi%%r", i))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		exp = color.Ti.Color(-1, i) + "hi" + color.Ti.Reset
		r = color.Highlight(fmt.Sprintf("%%h[bg%d]hi%%r", i))
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

var combinations = map[string]string{
	"%h[fgMaroon+bgNavy+bold+underline+fg23+bg235]":     color.Ti.Color(terminfo.ColorMaroon, terminfo.ColorNavy) + color.Ti.Bold + color.Ti.Underline + color.Ti.Color(23, 235),
	"%h[bgBlue+fgOlive+fgGreen+fg34+blink+dim+reverse]": color.Ti.Color(-1, terminfo.ColorBlue) + color.Ti.Color(terminfo.ColorOlive, -1) + color.Ti.Color(terminfo.ColorGreen, -1) + color.Ti.Color(34, -1) + color.Ti.Blink + color.Ti.Dim + color.Ti.Reverse,
}

func TestCombinations(t *testing.T) {
	for k, v := range combinations {
		if r := color.Highlight(k); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

var highlightEdgeCases = map[string]string{
	"%h[fgGray+%h[fgBlue]": color.Ti.Color(terminfo.ColorGray, -1) + color.ErrInvalid,
	"%h[":                  color.ErrInvalid,
	"%h{":                  color.ErrInvalid,
	"%h[]":                 color.ErrMissing,
	"%%h[fgRed]":           "%%h[fgRed]",
	"%[bg232]":             "%[bg232]",
	"%h[fg132":             color.Ti.Color(132, -1) + color.ErrInvalid,
	"%h[fgFuchsia[]":       color.Ti.Color(terminfo.ColorFuchsia, -1) + color.ErrInvalid,
	"%h[fgGreen+lold[]":    color.Ti.Color(terminfo.ColorGreen, -1) + color.ErrBadAttr,
	"%h[fgOlive+%#bgBlue]": color.Ti.Color(terminfo.ColorOlive, -1) + color.ErrInvalid,
	"%h][fgRed+%#bgBlue]":  color.ErrInvalid,
	"%h[fgRed+":            color.Ti.Color(terminfo.ColorRed, -1) + color.ErrInvalid,
	"%%h%h[fgRed]%%":       "%%h\x1b[91m%%",
	"%h[dsadadssadas]":     color.ErrBadAttr,
	"%":                    "%",
	"%h[fgsadas]":          color.ErrBadAttr,
	"%h[fgAqua+%h[bgBlue]": color.Ti.Color(terminfo.ColorAqua, -1) + color.ErrInvalid,
	"lmaokai":              "lmaokai",
	"%h[fgMaroon]%h[]":     color.Ti.Color(terminfo.ColorMaroon, -1) + color.ErrMissing,
	"%h[bgGjo]%h[bgGreen]": color.ErrBadAttr,
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

var result string

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
