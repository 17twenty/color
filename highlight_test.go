package color_test

import (
	"os"
	"testing"

	"github.com/gdamore/tcell"
	"github.com/nhooyr/color"
)

var ti, _ = tcell.LookupTerminfo(os.Getenv("TERM"))

// func TestAttributes(t *testing.T) {
// 	for k, v := range color.Colors {
// 		s := fmt.Sprintf("%%h[%s]hi%%r", k)
// 		exp := fmt.Sprintf("\x1b[%smhi\x1b[0m", v[1:])
// 		if r := color.Highlight(s, ti); r != exp {
// 			t.Errorf("Expected %q but result was %q", exp, r)
// 		}
// 	}
// }
//
// func TestColor256(t *testing.T) {
// 	for i := 0; i < 256; i++ {
// 		s := fmt.Sprintf("%%h[fg%d]hi%%r", i)
// 		exp := fmt.Sprintf("\x1b[38;5;%dmhi\x1b[0m", i)
// 		if r := color.Highlight(s, ti); r != exp {
// 			t.Errorf("Expected %q but result was %q", exp, r)
// 		}
// 		s = fmt.Sprintf("%%h[bg%d]hi%%r", i)
// 		exp = fmt.Sprintf("\x1b[48;5;%dmhi\x1b[0m", i)
// 		if r := color.Highlight(s, ti); r != exp {
// 			t.Errorf("Expected %q but result was %q", exp, r)
// 		}
// 	}
// }
//
// var combinations = map[string]string{
// 	"%h[fgRed+bgBlue+bold+underline+fg23+bg235]":   "\x1b[31;44;1;4;38;5;23;48;5;235m",
// 	"%h[bgBrightBlue+fgYellow+fgGreen+fg34+blink]": "\x1b[104;33;32;38;5;34;5m",
// }
//
// func TestCombinations(t *testing.T) {
// 	for k, v := range combinations {
// 		if r := color.Highlight(k, ti); r != v {
// 			t.Errorf("Expected %q but result was %q", v, r)
// 		}
// 	}
// }
//
// var highlightEdgeCases = map[string]string{
// 	"%h[fgRed+%h[fgBlue]": "%%!h(INVALID)",
// 	"%h[":                 "%%!h(INVALID)",
// 	"%h{":                 "%%!h(INVALID)",
// 	"%h[]":                "%%!h(MISSING)",
// 	"%%h[fgRed]":          "%%h[fgRed]",
// 	"%[bg232]":            "%[bg232]",
// 	"%h[fg132":            "%%!h(INVALID)",
// 	"%h[fgRed[]":          "%%!h(INVALID)",
// 	"%h[fgRed+lold[]":     "%%!h(BADATTR)",
// 	"%h[fgRed+%#bgBlue]":  "%%!h(INVALID)",
// 	"%h][fgRed+%#bgBlue]": "%%!h(INVALID)",
// 	"%h[fgRed+":           "%%!h(INVALID)",
// 	"%%h%h[fgRed]%%":      "%%h\x1b[31m%%",
// 	"%h[dsadadssadas]":    "%%!h(BADATTR)",
// 	"%":                   "%",
// 	"%h[fgsadas]":         "%%!h(BADATTR)",
// 	"%h[fgRed+%h[bgBlue]": "%%!h(INVALID)",
// 	"lmaokai":             "lmaokai",
// }
//
// func TestHighlightEdgeCases(t *testing.T) {
// 	for k, v := range highlightEdgeCases {
// 		if r := color.Highlight(k, ti); r != v {
// 			t.Errorf("Expected %q but result was %q", v, r)
// 		}
// 	}
// }
//
// var stripEdgeCases = map[string]string{
// 	"%h[fgRed]%smao%r": "%smao",
// 	"%":                "%",
// 	"%c":               "%c",
// }
//
// func TestStripEdgeCases(t *testing.T) {
// 	for k, v := range stripEdgeCases {
// 		if r := color.Strip(k, ti); r != v {
// 			t.Errorf("Expected %q but result was %q", v, r)
// 		}
// 	}
// }

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
