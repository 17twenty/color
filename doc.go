/*
Package color wraps fmt.Printf with verbs for producing colored output.

Highlight verbs:

	%h[attr...]	replaced with a SGR code that sets all of the attributes in []
			multiple attributes are + separated
	%r		an abbreviation for %h[reset]

Multiple highlight verbs do not reset preceding verbs, they add onto them.
For example, if you set the foreground to green and background to yellow in the first verb, then set the background to red in the second, any text following the second will have a green foreground and a red background.
This also applies across calls, the attributes are never reset unless explicitly requested.

Errors:

If an error occurs, one of the following strings will replace the position of the highlight verb.

	%!h(INVALID)	invalid character in the highlight verb
	%!h(MISSING)	no attributes in the highlight verb
	%!h(BADATTR)	unknown attribute in the highlight verb

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
	%h[fgBrightBlack]
	%h[fgBrightRed]
	%h[fgBrightGreen]
	%h[fgBrightYellow]
	%h[fgBrightBlue]
	%h[fgBrightMagenta]
	%h[fgBrightCyan]
	%h[fgBrightWhite]
	%h[bgBrightBlack]
	%h[bgBrightRed]
	%h[bgBrightGreen]
	%h[bgBrightYellow]
	%h[bgBrightBlue]
	%h[bgBrightMagenta]
	%h[bgBrightCyan]
	%h[bgBrightWhite]

256 Colors:
	%h[fgx]
	%h[bgx]
Where x is any number from 0-255

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
