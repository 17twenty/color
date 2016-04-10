/*
Package color wraps fmt.Printf with verbs for producing colored output.

Printing


Highlight Verbs:

	%h[attr...]	replaced with a SGR code that sets all of the attributes in []
			multiple attributes are + separated
	%r		an abbreviation for %h[reset]

Errors:

If an error occurs, the generated string will contain a description of the problem, as in these examples.

	No attributes in the highlight verb:
		Printf("%h[]%s", "hi"):			%!h(MISSING)
	Invalid character in the highlight verb:
		Printf("%h[%&:*!]%s", "hi"):		%!h(INVALID)
	Unknown attribute in the highlight verb:
		Printf("%h[fgOrange]%s", "hi"):	%!h(BADATTR)

Everything else is handled by the fmt package. You should read its documentation.

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

Other:
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

See http://goo.gl/LRLA7o for a more in depth explanation of the attributes. Scroll down till you see the SGR section.

See http://goo.gl/fvtHLs and according to the above document, ISO-8613-3, for more information on 256 colors.
*/
package color
