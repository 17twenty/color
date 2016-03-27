/*
	Package color extends fmt.Printf with verbs for terminal color highlighting. All it does is replace the verbs with the appropriate terminal escape sequence.

	Highlight verbs:

		%h[attrs]		uses the attrs to highlight the following text
		%r			an abbreviation for %h[reset]

	attrs is a + separated list of Colors (e.g. fgRed) or Attributes (e.g. bold).

	Multiple highlight verbs do not reset preceeding verbs, they add onto them.
	For example, if you set the foreground to green in the first verb, then set the background to red in the second, any text following the second will have a green foreground and a red background.

	Syntax Reference

	16 Foreground Colors:
		%h[fgBlack]
		%h[fgRed]
		%h[fgGreen]
		%h[fgYellow]
		%h[fgBlue]
		%h[fgMagenta]
		%h[fgCyan]
		%h[fgWhite]
		%h[fgDefault]
		%h[fgBrighBlack]
		%h[fgBrightRed]
		%h[fgBrightGreen]
		%h[fgBrightYellow]
		%h[fgBrightBlue]
		%h[fgBrightMagenta]
		%h[fgBrightCyan]
		%h[fgBrightWhite]

	16 Background Colors:
		%h[bgBlack]
		%h[bgRed]
		%h[bgGreen]
		%h[bgYellow]
		%h[bgBlue]
		%h[bgMagenta]
		%h[bgCyan]
		%h[bgWhite]
		%h[bgDefault]
		%h[bgBrighBlack]
		%h[bgBrightRed]
		%h[bgBrightGreen]
		%h[bgBrightYellow]
		%h[bgBrightBlue]
		%h[bgBrightMagenta]
		%h[bgBrightCyan]
		%h[bgBrightWhite]

	256 Colors:
		%h[fgxxx]
		%h[bgxxx]
	Where xxx is any number from 0-255

	Attributes:
		%h[reset] or the %r verb
		%h[bold]
		%h[faint]
		%h[italic]
		%h[underline]
		%h[blink]
		%h[inverse]
		%h[invisible]
		%h[crossedOut]
		%h[doubleUnderline]
		%h[normal]
		%h[notItalic]
		%h[notUnderlined]
		%h[steady]
		%h[positive]
		%h[visible]
		%h[notCrossedOut]

	Mixing:
		%h[fgBlue+bgRed+bold]

	See http://goo.gl/LRLA7o for an explanation of these colors/attributes. Scroll down to the SGR section.
*/
package color
