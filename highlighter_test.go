package color

import "testing"

func BenchmarkHighlight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Highlight("swag %h[fgRed]sdsa%r%h[fgBlue+bold+underline]swagw%r%h[bgRed+fgGreen+bold]lmaokai%r%h[fg101+bg34]dsaasdas%r")
	}
}
