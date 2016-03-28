/*
	Package color extends fmt.Printf with verbs for terminal color highlighting.

	Highlight verbs:

		%h[attr...]		replaces itself with a SGR code that sets all of the attributes in []
					multiple attributes are + separated
		%r			an abbreviation for %h[reset]

	Multiple highlight verbs do not reset preceeding verbs, they add onto them.
	For example, if you set the foreground to green in the first verb, then set the background to red in the second, any text following the second will have a green foreground and a red background.

	Attributes Reference

	Standard Colors:
		%h[fgBlack]
		%h[fgRed]
		%h[fgGreen]
		%h[fgYellow]
		%h[fgBlue]
		%h[fgMagenta]
		%h[fgCyan]
		%h[fgWhite]
		%h[fgDefault]
		%h[bgBlack]
		%h[bgRed]
		%h[bgGreen]
		%h[bgYellow]
		%h[bgBlue]
		%h[bgMagenta]
		%h[bgCyan]
		%h[bgWhite]
		%h[bgDefault]

	Bright Colors:
		%h[bgBrighBlack]
		%h[bgBrightRed]
		%h[bgBrightGreen]
		%h[bgBrightYellow]
		%h[bgBrightBlue]
		%h[bgBrightMagenta]
		%h[bgBrightCyan]
		%h[bgBrightWhite]
		%h[fgBrighBlack]
		%h[fgBrightRed]
		%h[fgBrightGreen]
		%h[fgBrightYellow]
		%h[fgBrightBlue]
		%h[fgBrightMagenta]
		%h[fgBrightCyan]
		%h[fgBrightWhite]

	256 Colors:
		%h[fgxxx]
		%h[bgxxx]
	Where xxx is any number from 0-255

	Others:
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

	See http://goo.gl/LRLA7o for a more in depth explanation of the attributes. Scroll down till you see the SGR section
*/
package color
