/*
Package color adds verbs to fmt.Printf for producing colored output.

Printing

Verbs:

	%h[attr...]	replaced with a SGR code that sets all of the attributes in []
			multiple attributes are + separated
	%r		an abbreviation for %h[reset]

Errors:

If an error occurs, the generated string will contain a description of the problem, as in these examples.

	Invalid character in the highlight verb:
		Printf("%h(fgMaroon)%s", "hi"):		%!h(INVALID)
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
	%h[xgMaroon]
	%h[xgGreen]
	%h[xgOlive]
	%h[xgNavy]
	%h[xgPurple]
	%h[xgTeal]
	%h[xgSilver]
	%h[xgGray]
	%h[xgRed]
	%h[xgLime]
	%h[xgYellow]
	%h[xgBlue]
	%h[xgFuchsia]
	%h[xgAqua]
	%h[xgWhite]

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

See http://jonasjacek.github.io/colors/ for a reference of the colors.
*/
package color
