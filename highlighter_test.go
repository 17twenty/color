package color

import "testing"

func BenchmarkHighlight(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// standard colors
		Highlight(`
		%h[fgBlack]hi%r
		%h[fgRed]hi%r
		%h[bgGreen]hi%r
		%h[bgYellow]hi%r`)
		Highlight(`
		%h[fgBrightBlue]hi%r
		%h[fgBrightMagenta]hi%r
		%h[bgBrightCyan]hi%r
		%h[bgBrightWhite]hi%r`)
		Highlight(`
		%h[fg22]hi%r
		%h[fg233]hi%r
		%h[bg3]hi%r
		%h[bg102]hi%r`)
		Highlight(`
		%h[bold]hi%r
		%h[underline]hi%r
		%h[italic]hi%r
		%h[blink]hi%r`)
		Highlight(`
		%h[fgRed+bgBrightBlack+bold+underline]hi%r
		%h[fg32+bg123+bold+underline]hi%r`)
	}
}
