package color

import (
	"fmt"
	"testing"
)

func TestAttributes(t *testing.T) {
	for k, v := range attrs {
		s := fmt.Sprintf("%%h[%s]hi%%r", k)
		exp := fmt.Sprintf("\x1b[%smhi\x1b[0m", v[1:])
		if r := Highlight(s); r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

func TestColor256(t *testing.T) {
	for i := 0; i < 256; i++ {
		s := fmt.Sprintf("%%h[fg%d]hi%%r", i)
		exp := fmt.Sprintf("\x1b[38;5;%dmhi\x1b[0m", i)
		if r := Highlight(s); r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		s = fmt.Sprintf("%%h[bg%d]hi%%r", i)
		exp = fmt.Sprintf("\x1b[48;5;%dmhi\x1b[0m", i)
		if r := Highlight(s); r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

var combinations = map[string]string{
	"%h[fgRed+bgBlue+bold+underline+fg23+bg235]":   "\x1b[31;44;1;4;38;5;23;48;5;235m",
	"%h[bgBrightBlue+fgYellow+fgGreen+fg34+blink]": "\x1b[104;33;32;38;5;34;5m",
}

func TestCombinations(t *testing.T) {
	for k, v := range combinations {
		if r := Highlight(k); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

var edgeCases = map[string]string{
	"%h[dsadadssadas]":    "%h[dsadadssadas]",
	"%h[":                 "%h[",
	"%h{":                 "%h{",
	"%%h[fgRed]":          "%%h[fgRed]",
	"%[fgRed]":            "%[fgRed]",
	"%h[fgRed":            "%h[fgRed",
	"%h[fgRed[]":          "%h[fgRed[]",
	"%h[fgRed+lold[]":     "%h[fgRed+lold[]",
	"%h[fgRed+%#bgBlue]":  "%h[fgRed+%#bgBlue]",
	"%h][fgRed+%#bgBlue]": "%h][fgRed+%#bgBlue]",
	"%h[fgRed+":           "%h[fgRed+",
	"%%h%h[fgRed]%%":      "%%h\x1b[31m%%",
	"%h[fgsadas]":         "%h[fgsadas]",
	"%h[fgRed+%h[bgBlue]": "%h[fgRed+\x1b[44m",
}

// func TestEdgeCases(t *testing.T) {
// 	for k, v := range edgeCases {
// 		if r := Highlight(k); r != v {
// 			t.Errorf("Expected %q but result was %q", v, r)
// 		}
// 	}
// }

func BenchmarkHighlight(b *testing.B) {
	s := `%h[fgBlack]hi%r
	%h[fgRed]hi%r
	%h[bgGreen]hi%r
	%h[bgYellow]hi%r
	%h[fgBrightBlue]hi%r
	%h[fgBrightMagenta]hi%r
	%h[bgBrightCyan]hi%r
	%h[bgBrightWhite]hi%r
	%h[bold]hi%r
	%h[underline]hi%r
	%h[italic]hi%r
	%h[blink]hi%r
	%h[fg22]hi%r
	%h[fg233]hi%r
	%h[bg3]hi%r
	%h[bg102]hi%r
	%h[fgWhite+bgBrightCyan+bold+underline+fg32+bg69]hi%r
	%h[fg32+bg123+bold+underline+bgBlue+fgBrighGreen+bgBrightWhite]hi%r`
	for i := 0; i < b.N; i++ {
		Highlight(s)
	}
}
