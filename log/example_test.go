package log_test

import (
	"github.com/nhooyr/color"
	"github.com/nhooyr/color/log"
)

func Example() {
	redFormat := color.Prepare("%h[fgMaroon]%s%r\n")

	// If os.Stderr is a terminal, this will print in color.
	// Otherwise it will be a normal "hi".
	log.Printfp(redFormat, "hi")

	// normal "hi", the highlight verbs are ignored.
	log.SetColor(false)
	log.Printfp(redFormat, "hi")

	// "hi" with a red foreground.
	log.SetColor(true)
	log.Fatalfp(redFormat, "hi")
}
