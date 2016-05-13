package log_test

import (
	"github.com/nhooyr/color"
	"github.com/nhooyr/color/log"
)

func Example() {
	f := color.Prepare("%h[fgRed]%s%r\n")

	// If os.Stderr is a terminal, this will print in color.
	// Otherwise it will be a normal "foo".
	log.Printf(f, "foo")

	// Normal "bar", the highlight verbs are ignored.
	log.SetColor(false)
	log.Printf(f, "bar")

	// "foo" with a red foreground.
	log.SetColor(true)
	log.Fatalf(f, "foo")
}
