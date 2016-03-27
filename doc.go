/*
Package color extends fmt.Printf with verbs for terminal color highlighting. All it does is replace the verbs with the appropiate terminal escape sequence.

Highlight verbs

	%h[attrs] 	uses the attrs to highlight the following text
				multiple attrs are separated by +
				subsquent highlighting verbs do not first reset, they just append
	%r			an abbreviation for %h[reset]

Standard Colors
	// red "panic:" and then normal "rip"
	color.Printf("%h[fgRed]panic:%r rip\n")

	// "panic:" with brightRed background and then normal "rip"
	color.Printf("%h[bgBrightRed]panic:%r rip\n")

Attributes
	// red "panic:" and then "rip" with a cyan background
	color.Printf("%h[fg1]panic:%r %h[bg6]rip\n")

256 Colors
	// red "panic:" and then "rip" with a cyan background
	color.Printf("%h[fg1]panic:%r %h[bg6]rip\n")

Mixing colors and attributes
	// bolded green "panic:" and then underlined "rip" with bright black background
	color.Printf("%h[fgGreen+bold]panic:%r %h[bg8+underline]rip%r\n")

How does reset behave?
	// bolded green "panic:" and then bolded green "rip" with bright black background
	color.Printf("%h[fgGreen+bold]panic: %h[bg8]rip\n")
	// bolded green "hi" with bright black background
	fmt.Printf("hi")
	// finally resets the highlighting
	color.Printf("%r")
*/
package color
