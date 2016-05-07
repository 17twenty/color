package color_test

import (
	"fmt"
	"os"

	"github.com/nhooyr/color"
	"github.com/nhooyr/color/log"
)

func Example_attributes() {
	// "panic:" with a maroon foreground then normal "rip".
	color.Printf("%h[fgMaroon]panic:%r %s\n", "rip")

	// "panic:" with a red background then normal "rip".
	color.Printf("%h[bgRed]panic:%r %s\n", "rip")

	// Bold "panic:" then normal "rip".
	color.Printf("%h[bold]panic:%r %s\n", "rip")

	// Underlined "panic:" with then normal "rip".
	color.Printf("%h[underline]panic:%r %s\n", "rip")

	// "panic:" using color 83 as the foreground then normal "rip".
	color.Printf("%h[fg83]panic:%r %s\n", "rip")

	// "panic:" using color 158 as the background then normal "rip".
	color.Printf("%h[bg158]panic:%r %s\n", "rip")
}

func Example_mixing() {
	// Bolded "panic:" with a green foreground then normal "rip".
	color.Printf("%h[fgGreen+bold]panic:%r %s\n", "rip")

	// Underlined "panic:" with a gray background then normal "rip".
	color.Printf("%h[bg8+underline]panic:%r %s\n", "rip")
}

func ExamplePrepare() {
	// Prepare only processes the highlight verbs in the string,
	// letting you print it repeatedly with performance.
	panicFormat := color.Prepare("%h[fgMaroon+bold]panic:%r %s\n")

	// If os.Stdout is a terminal, each will print a bolded "panic:" in red foreground
	// and some normal text after. Otherwise each will print normally.
	color.Printfp(panicFormat, "rip")
	color.Printfp(panicFormat, "yippie")
	color.Printfp(panicFormat, "dsda")
}

func ExamplePrinter() {
	redFormat := color.Prepare("%h[fgMaroon]%s%r\n")

	// If os.Stderr is a terminal, this will print in color.
	// Otherwise it will be a normal "hi".
	p = color.New(os.Stderr, color.IsTerminal(os.Stderr))
	p.Printfp(redFormat, "hi")

	// "hi" with red foreground.
	p := color.New(os.Stderr, true)
	p.Printfp(redFormat, "hi")

	// normal "hi", the highlight verbs are ignored.
	p = color.New(os.Stderr, false)
	p.Printfp(redFormat, "hi")
}

func Example_reset() {
	// "rip" will be printed with a navy foreground and gray background
	// because we never reset the highlighting after "panic:". The navy foreground is
	// carried on from "panic:".
	color.Printf("%h[fgNavy+bgGray]panic: %h[bg8]%s\n", "rip")

	// The attributes carry onto anything written to the terminal until reset.
	// This prints "rip" in the same attributes as above.
	fmt.Println("rip")

	// Resets the highlighting and then prints "hello" normally.
	color.Printf("%r%s", "hello")
}
