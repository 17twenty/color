package color_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/nhooyr/color"
)

var ti, _ = tcell.LookupTerminfo(os.Getenv("TERM"))

func TestAttributes(t *testing.T) {
	for k, v := range color.Colors {
		exp := ti.TColor(v, -1) + "hi" + ti.AttrOff
		r := color.Highlight(fmt.Sprintf("%%h[fg%s]hi%%r", k), ti)
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		exp = ti.TColor(-1, v) + "hi" + ti.AttrOff
		r = color.Highlight(fmt.Sprintf("%%h[bg%s]hi%%r", k), ti)
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

func TestColor256(t *testing.T) {
	for i := tcell.Color(0); i < 256; i++ {
		exp := ti.TColor(i, -1) + "hi" + ti.AttrOff
		r := color.Highlight(fmt.Sprintf("%%h[fg%d]hi%%r", i), ti)
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		exp = ti.TColor(-1, i) + "hi" + ti.AttrOff
		r = color.Highlight(fmt.Sprintf("%%h[bg%d]hi%%r", i), ti)
		if r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

var combinations = map[string]string{
	"%h[fgMaroon+bgNavy+bold+underline+fg23+bg235]":     ti.TColor(tcell.ColorMaroon, tcell.ColorNavy) + ti.Bold + ti.Underline + ti.TColor(23, 235),
	"%h[bgBlue+fgOlive+fgGreen+fg34+blink+dim+reverse]": ti.TColor(-1, tcell.ColorBlue) + ti.TColor(tcell.ColorOlive, -1) + ti.TColor(tcell.ColorGreen, -1) + ti.TColor(34, -1) + ti.Blink + ti.Dim + ti.Reverse,
}

func TestCombinations(t *testing.T) {
	for k, v := range combinations {
		if r := color.Highlight(k, ti); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

var highlightEdgeCases = map[string]string{
	"%h[fgGray+%h[fgBlue]": ti.TColor(tcell.ColorGray, -1) + color.ErrInvalid,
	"%h[":                  color.ErrInvalid,
	"%h{":                  color.ErrInvalid,
	"%h[]":                 color.ErrMissing,
	"%%h[fgRed]":           "%%h[fgRed]",
	"%[bg232]":             "%[bg232]",
	"%h[fg132":             ti.TColor(132, -1) + color.ErrInvalid,
	"%h[fgFuchsia[]":       ti.TColor(tcell.ColorFuchsia, -1) + color.ErrInvalid,
	"%h[fgGreen+lold[]":    ti.TColor(tcell.ColorGreen, -1) + color.ErrBadAttr,
	"%h[fgOlive+%#bgBlue]": ti.TColor(tcell.ColorOlive, -1) + color.ErrInvalid,
	"%h][fgRed+%#bgBlue]":  color.ErrInvalid,
	"%h[fgRed+":            ti.TColor(tcell.ColorRed, -1) + color.ErrInvalid,
	"%%h%h[fgRed]%%":       "%%h\x1b[91m%%",
	"%h[dsadadssadas]":     color.ErrBadAttr,
	"%":                    "%",
	"%h[fgsadas]":          color.ErrBadAttr,
	"%h[fgAqua+%h[bgBlue]": ti.TColor(tcell.ColorAqua, -1) + color.ErrInvalid,
	"lmaokai":              "lmaokai",
	"%h[fgMaroon]%h[]":     ti.TColor(tcell.ColorMaroon, -1) + color.ErrMissing,
	"%h[bgGjo]%h[bgGreen]": color.ErrBadAttr,
}

func TestHighlightEdgeCases(t *testing.T) {
	for k, v := range highlightEdgeCases {
		if r := color.Highlight(k, ti); r != v {
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
		if r := color.Strip(k, ti); r != v {
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
		r = color.Highlight(s, ti)
	}
	result = r
}

func BenchmarkStrip(b *testing.B) {
	var r string
	for i := 0; i < b.N; i++ {
		r = color.Strip(s, ti)
	}
	result = r
}
