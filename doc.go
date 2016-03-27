/*
	Package color extends fmt.Printf with verbs for terminal color highlighting. All it does is replace the verbs with the appropiate terminal escape sequence.

	Highlight verbs:

		%h[attrs]		uses the attrs to highlight the following text
		%r			an abbreviation for %h[reset]
	
	attrs is a plus sign separated list of Colors (fgRed) or Attributes (bold).

	Multiple highlight verbs do not reset preceeding verbs, they add onto them.
	For example, if you set the foreground to green in the first verb, then set the background to red in the second, any text following the second will have a green foreground and a red background.
	
	Standard Colors:
		// red "panic:" and then normal "rip"
		color.Printf("%h[fgRed]panic:%r rip\n")

		// "panic:" with brightRed background and then normal "rip"
		color.Printf("%h[bgBrightRed]panic:%r rip\n")

	Attributes:
		// red "panic:" and then "rip" with a cyan background
		color.Printf("%h[fg1]panic:%r %h[bg6]rip\n")

	256 Colors:
		// red "panic:" and then "rip" with a cyan background
		color.Printf("%h[fg1]panic:%r %h[bg6]rip\n")

	Mixing colors and attributes:
		// bolded green "panic:" and then underlined "rip" with bright black background
		color.Printf("%h[fgGreen+bold]panic:%r %h[bg8+underline]rip%r\n")

	Reset behavior:
		// bolded green "panic:" and then bolded green "rip" with bright black background
		color.Printf("%h[fgGreen+bold]panic: %h[bg8]rip\n")
		// bolded green "hi" with bright black background
		fmt.Printf("hi")
		// finally resets the highlighting
		color.Printf("%r")
*/
package color
