package color_test

import (
	"fmt"
	"log"
	"os"

	"github.com/nhooyr/color"
)

func Example_attributes() {
	// "panic:" with a maroon foreground then normal "rip".
	color.Printfh("%h[fgMaroon]panic:%r %s\n", "rip")

	// "panic:" with a red background then normal "rip".
	color.Printfh("%h[bgRed]panic:%r %s\n", "rip")

	// Bold "panic:" then normal "rip".
	color.Printfh("%h[bold]panic:%r %s\n", "rip")

	// Underlined "panic:" with then normal "rip".
	color.Printfh("%h[underline]panic:%r %s\n", "rip")

	// "panic:" using color 83 as the foreground then normal "rip".
	color.Printfh("%h[fg83]panic:%r %s\n", "rip")

	// "panic:" using color 158 as the background then normal "rip".
	color.Printfh("%h[bg158]panic:%r %s\n", "rip")
}

func Example_mixing() {
	// Bolded "panic:" with a green foreground then normal "rip".
	color.Printfh("%h[fgGreen+bold]panic:%r %s\n", "rip")

	// Underlined "panic:" with a gray background then normal "rip".
	color.Printfh("%h[bg8+underline]panic:%r %s\n", "rip")
}

func ExamplePrepare() {
	// Prepare only processes the highlight verbs in the string,
	// letting you print it repeatedly with performance.
	panicFormat := color.Prepare("%h[fgMaroon+bold]panic:%r %s\n")

	// Each prints a bolded "panic:" in red foreground and some normal text after.
	// Notice that fmt.Printf is used, this works because only the highlight verbs
	// were processed above, the %s verb was not.
	fmt.Printf(panicFormat, "rip")
	fmt.Printf(panicFormat, "yippie")
	fmt.Printf(panicFormat, "dsda")
}

func ExamplePrinter() {
	// "hi" with red foreground.
	p := color.NewPrinter(os.Stderr, true)
	// See the Prepare example for an explanation of this.
	redFormat := p.Prepare("%h[fgMaroon]%s%r\n")
	p.Printf(redFormat, "hi")

	// normal "hi", the highlight verbs are ignored.
	p = color.NewPrinter(os.Stderr, false)
	p.Printfh("%h[fgMaroon]%s%r\n", "hi")

	// If os.Stderr is a terminal, this will print in color.
	// Otherwise it will be a normal "hi".
	p = color.NewPrinter(os.Stderr, color.IsTerminal(os.Stderr))
	p.Printfh("%h[fgMaroon]%s%r\n", "hi")
}

func ExampleLogger() {
	// "hi" with a red foreground.
	l := color.NewLogger(os.Stderr, "%h[bold]color:%r ", log.LstdFlags, true)
	// See the Prepare example for an explanation of this.
	redFormat := l.Prepare("%h[fgMaroon]%s%r\n")
	l.Printf(redFormat, "hi")

	// normal "hi", the highlight verbs are ignored.
	l = color.NewLogger(os.Stderr, "%h[bold]color:%r ", log.LstdFlags, false)
	l.Printfh("%h[fgMaroon]%s%r", "hi")

	// If os.Stderr is a terminal, this will print in color.
	// Otherwise it will be a normal "hi".
	l = color.NewLogger(os.Stderr, "%h[bold]color:%r ", log.LstdFlags, color.IsTerminal(os.Stderr))
	l.Fatalf("%h[fgMaroon]%s%r", "hi")
}

func Example_reset() {
	// "rip" will be printed with a navy foreground and gray background
	// because we never reset the highlighting after "panic:". The navy foreground is
	// carried on from "panic:".
	color.Printfh("%h[fgNavy+bgGray]panic: %h[bg8]%s\n", "rip")

	// The attributes carry onto anything written to the terminal until reset.
	// This prints "rip" in the same attributes as above.
	fmt.Println("rip")

	// Resets the highlighting and then prints "hello" normally.
	color.Printfh("%r%s", "hello")
}
