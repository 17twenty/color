/*
Package color extends fmt.Printf with verbs for producing colored output.

Printing

Verbs:

	%h[attr...]	replaced with a SGR code that sets all of the attributes in []
			multiple attributes are + separated
	%r		an abbreviation for %h[reset]
	%a		used by Format's Insert methods to combine Formats

Before printing, Prepare must be called with the format string. It will return a Format structure that represents the format string with the highlight verbs fully parsed. Why?

While this package is heavily optimized, processing the highlighting verbs is still expensive. Thus, it makes more sense to process the verbs once and then store the results into a Format structure.

It holds the colored and stripped versions of the base format string. In the colored string, the highlight verbs are replaced with their control sequences. In contrast, the highlight verbs are completely removed in the stripped string. Why store both? If color output is enabled, the colored string is used, but if color output is disabled, then the stripped string is used.

There are two methods to print these Format structures. One, the Printf like functions take a Format structure that they will appropriately expand and use as the format string. Two, every Print like function's variadic arguments can contain Format structures; they will be expanded to their appropriate strings.

See Example (Printing) for an example of both methods.

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
