/*
Package color extends fmt.Printf with verbs for producing colored output.

Printing

Verbs:

	%h[attr...]	replaced with a SGR code that sets all of the attributes in []
			multiple attributes are + separated
	%r		an abbreviation for %h[reset]

Errors:

If an error occurs, the generated string will contain a description of the problem, as in these examples.

	Invalid character in the highlight verb:
		Printf("%h(fgRed)%s", "hi"):		%!h(INVALID)
	No attributes in the highlight verb:
		Printf("%h[]%s", "hi"):			%!h(MISSING)
	Unknown attribute in the highlight verb:
		Printf("%h[fgGdsds]%s", "hi"):		%!h(BADATTR)
	String ended before the verb:
		Printf("%h[fg", "hi"):			%!h(SHORT)

Everything else is handled by the fmt package. You should read its documentation.

Attributes Reference

Named Colors:
	%h[xgBlack]
	%h[xgRed]
	%h[xgGreen]
	%h[xgYellow]
	%h[xgBlue]
	%h[xgMagenta]
	%h[xgCyan]
	%h[xgWhite]
	%h[xgBrightBlack]
	%h[xgBrightRed]
	%h[xgBrightGreen]
	%h[xgBrightYellow]
	%h[xgBrightBlue]
	%h[xgBrightMagenta]
	%h[xgBrightCyan]
	%h[xgBrightWhite]

	Where 'x' is either 'f' or 'b'.

256 Colors:
	%h[fgx]
	%h[bgx]

	Where x is any number from 0-255.

Modes:
	%h[reset] or the %r verb
	%h[bold]
	%h[underline]
	%h[reverse]
	%h[blink]
	%h[dim]

See http://goo.gl/LRLA7o for information on the attributes. Scroll down to the SGR section.

See http://goo.gl/fvtHLs and ISO-8613-3 (according to above document) for more information on 256 colors.
*/
package color
