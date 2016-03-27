/*
	Package color extends fmt.Printf with verbs for terminal color highlighting. All it does is replace the verbs with the appropriate terminal escape sequence.

	Highlight verbs:

		%h[attrs]		uses the attrs to highlight the following text
		%r			an abbreviation for %h[reset]

	attrs is a + separated list of Colors (e.g. fgRed) or Attributes (e.g. bold).

	Multiple highlight verbs do not reset preceeding verbs, they add onto them.
	For example, if you set the foreground to green in the first verb, then set the background to red in the second, any text following the second will have a green foreground and a red background.

	The syntax reference is here: https://github.com/nhooyr/color/blob/master/REFERENCE.md

	Examples

	Standard Colors:
		// "panic:" with a red foreground then normal "rip"
		color.Printf("%h[fgRed]panic:%r rip\n")

		// "panic:" with a brightRed background then normal "rip"
		color.Printf("%h[bgBrightRed]panic:%r rip\n")

	Attributes:
		// bold "panic:" then normal "rip"
		color.Printf("%h[bold]panic:%r rip\n")

		// underlined "panic:" with then normal "rip"
		color.Printf("%h[underline]panic:%r rip\n")

	256 Colors:
		// "panic:" with a green foreground then normal "rip"
		color.Printf("%h[fg2]panic:%r rip\n")

		// "panic:" with a bright green background then normal "rip"
		color.Printf("%h[bg10]panic:%r rip\n")

	Mixing Colors and Attributes:
		// bolded "panic:" with a green foreground then normal "rip"
		color.Printf("%h[fgGreen+bold]panic:%r rip\n")

		// underlined "panic:" with a bright black background then normal "rip"
		color.Printf("%h[bg8+underline]panic:%r rip\n")

	How does reset behave?
		// bolded "panic:" with a blue foreground
		// then bolded "rip" with a green foreground and bright black background
		color.Printf("%h[fgBlue+bold]panic: %h[bg8]rip\n")

		// bolded "hi" with a green foreground and bright black background
		fmt.Printf("hi")

		// finally resets the highlighting
		color.Printf("%r")

	log.Logger wrapper:
		logger := color.NewLogger(os.Stderr, "", 0)

		// prints hi in red
		logger.Printf("%h[fgRed]hi%r")

		// prints hi in red and then exits
		logger.Fatalf("%h[fgRed]hi%r")

	Reference

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
*/
package color
