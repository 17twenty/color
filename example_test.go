package color_test

import (
	"fmt"
	"os"

	"github.com/nhooyr/color"
)

func Example_attributes() {
	// "panic:" with a red foreground then normal "foo".
	color.Printf("%h[fgRed]panic:%r %s\n", "foo")

	// "panic:" with a red background then normal "bar".
	color.Printf("%h[bgRed]panic:%r %s\n", "bar")

	// Bold "panic:" then normal "foo".
	color.Printf("%h[bold]panic:%r %s\n", "foo")

	// Underlined "panic:" with then normal "bar".
	color.Printf("%h[underline]panic:%r %s\n", "bar")

	// "panic:" using color 83 as the foreground then normal "foo".
	color.Printf("%h[fg83]panic:%r %s\n", "foo")

	// "panic:" using color 158 as the background then normal "bar".
	color.Printf("%h[bg158]panic:%r %s\n", "bar")
}

func Example_mixing() {
	// Bolded "panic:" with a green foreground then normal "foo".
	color.Printf("%h[fgGreen+bold]panic:%r %s\n", "foo")

	// Underlined "panic:" with a bright black background then normal "bar".
	color.Printf("%h[bg8+underline]panic:%r %s\n", "bar")
}

func ExamplePrepare() {
	// Prepare only processes the highlight verbs in the string,
	// letting you print it repeatedly with performance.
	panicFormat := color.Prepare("%h[fgRed+bold]panic:%r %s\n")

	// Each will print "panic:" and some normal text after, but if standard output
	// is a terminal, "panic:" will be printed in bold with a red foreground.
	color.Printfp(panicFormat, "foo")
	color.Printfp(panicFormat, "bar")
	color.Printfp(panicFormat, "foo")

	staticMessage := color.Prepare("%h[fgBlue+bold]HELLO%r")

	// Each will print "HELLO", but if standard output is a terminal,
	// "HELLO" will be printed in bold with a blue foreground.
	color.Println(staticMessage)
	color.Println(staticMessage)
	color.Println(staticMessage)
}

func ExamplePrinter() {
	redFormat := color.Prepare("%h[fgRed]%s%r\n")

	// If standard error is a terminal, this will print in color.
	// Otherwise it will print a normal "bar".
	p := color.New(os.Stderr, color.IsTerminal(os.Stderr))
	p.Printfp(redFormat, "bar")

	// "foo" with red foreground.
	p = color.New(os.Stderr, true)
	p.Printfp(redFormat, "foo")

	// normal "bar", the highlight verbs are ignored.
	p = color.New(os.Stderr, false)
	p.Printfp(redFormat, "bar")
}

func Example_reset() {
	// "hello" will be printed with a black foreground and bright green background
	// because we never reset the highlighting after "panic:". The black foreground is
	// carried on from "panic:".
	color.Printf("%h[fgBlack+bgBrightRed]panic: %h[bgBrightGreen]%s", "hello")

	// The attributes carry onto anything written to the terminal until reset.
	// This prints "world" in the same attributes as above.
	fmt.Println("world")

	// Resets the highlighting and then prints "hello" normally.
	color.Printf("%r%s", "foo")
}
