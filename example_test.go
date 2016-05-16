package color_test

import (
	"fmt"
	"os"

	"github.com/nhooyr/color"
)

func Example_attributes() {
	// "panic:" with a red foreground then normal "foo".
	f := color.Prepare("%h[fgRed]panic:%r %s\n")
	color.Printf(f, "foo")

	// "panic:" with a red background then normal "bar".
	f = color.Prepare("%h[bgRed]panic:%r %s\n")
	color.Printf(f, "bar")

	// Bold "panic:" then normal "foo".
	f = color.Prepare("%h[bold]panic:%r %s\n")
	color.Printf(f, "foo")

	// Underlined "panic:" with then normal "bar".
	f = color.Prepare("%h[underline]panic:%r %s\n")
	color.Printf(f, "bar")

	// "panic:" using color 83 as the foreground then normal "foo".
	f = color.Prepare("%h[fg83]panic:%r %s\n")
	color.Printf(f, "foo")

	// "panic:" using color 158 as the background then normal "bar".
	f = color.Prepare("%h[bg158]panic:%r %s\n")
	color.Printf(f, "bar")
}

func Example_mixing() {
	// Bolded "panic:" with a green foreground then normal "foo".
	f := color.Prepare("%h[fgGreen+bold]panic:%r %s\n")
	color.Printf(f, "foo")

	// Underlined "panic:" with a bright black background then normal "bar".
	f = color.Prepare("%h[bg8+underline]panic:%r %s\n")
	color.Printf(f, "bar")
}

func ExamplePrinter() {
	f := color.Prepare("%h[fgRed]%s%r\n")

	// If standard error is a terminal, this will print in color.
	// Otherwise it will print a normal "bar".
	p := color.New(os.Stderr, color.IsTerminal(os.Stderr))
	p.Printf(f, "bar")

	// "foo" with red foreground.
	p = color.New(os.Stderr, true)
	p.Printf(f, "foo")

	// Normal "bar", the highlight verbs are ignored.
	p = color.New(os.Stderr, false)
	p.Printf(f, "bar")
}

func Example_reset() {
	// "hello" will be printed with a black foreground and bright green background
	// because we never reset the highlighting after "panic:". The black foreground is
	// carried on from "panic:".
	f := color.Prepare("%h[fgBlack+bgBrightRed]panic: %h[bgBrightGreen]%s")
	color.Printf(f, "hello")

	// The attributes carry onto anything written to the terminal until reset.
	// This prints "world" in the same attributes as above.
	fmt.Println("world")

	// Resets the highlighting and then prints "hello" normally.
	f = color.Prepare("%r%s")
	color.Printf(f, "foo")
}
