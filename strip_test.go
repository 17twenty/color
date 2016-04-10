package color

import "testing"

var stripEdgeCases = map[string]string{
	"%h[fgRed]%smao%r": "%smao",
	"%":                "%",
	"%c":               "%c",
}

func TestStripEdgeCases(t *testing.T) {
	for k, v := range stripEdgeCases {
		if r := Sstripf(k); r != v {
			t.Errorf("Expected %q but result was %q", v, r)
		}
	}
}

func BenchmarkSstripf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sstripf(s)
	}
}
