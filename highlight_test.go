package color

import (
	"fmt"
	"testing"
)

func TestAttributes(t *testing.T) {
	for k, v := range attrs {
		s := fmt.Sprintf("%%h[%s]hi%%r", k)
		exp := fmt.Sprintf("\x1b[%smhi\x1b[0m", v[1:])
		if r := Shighlightf(s); r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
	}
}

func TestColor256(t *testing.T) {
	for i := 0; i < 256; i++ {
		s := fmt.Sprintf("%%h[fg%d]hi%%r", i)
		exp := fmt.Sprintf("\x1b[38;5;%dmhi\x1b[0m", i)
		if r := Shighlightf(s); r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		}
		s = fmt.Sprintf("%%h[bg%d]hi%%r", i)
		exp = fmt.Sprintf("\x1b[48;5;%dmhi\x1b[0m", i)
		if r := Shighlightf(s); r != exp {
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
		if r := Shighlightf(k); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

var highlightEdgeCases = map[string]string{
	"%h[fgRed+%h[fgBlue]": "%%!h(INVALID)",
	"%h[":                 "%%!h(INVALID)",
	"%h{":                 "%%!h(INVALID)",
	"%h[]":                "%%!h(MISSING)",
	"%%h[fgRed]":          "%%h[fgRed]",
	"%[bg232]":            "%[bg232]",
	"%h[fg132":            "%%!h(INVALID)",
	"%h[fgRed[]":          "%%!h(INVALID)",
	"%h[fgRed+lold[]":     "%%!h(BADATTR)",
	"%h[fgRed+%#bgBlue]":  "%%!h(INVALID)",
	"%h][fgRed+%#bgBlue]": "%%!h(INVALID)",
	"%h[fgRed+":           "%%!h(INVALID)",
	"%%h%h[fgRed]%%":      "%%h\x1b[31m%%",
	"%h[dsadadssadas]":    "%%!h(BADATTR)",
	"%":                   "%",
	"%h[fgsadas]":         "%%!h(BADATTR)",
	"%h[fgRed+%h[bgBlue]": "%%!h(INVALID)",
	"lmaokai":             "lmaokai",
}

func TestEdgeCases(t *testing.T) {
	for k, v := range highlightEdgeCases {
		if r := Shighlightf(k); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

func BenchmarkShighlightf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Shighlightf(s)
	}
}