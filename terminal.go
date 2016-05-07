package color

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// IsTerminal returns true if f is a terminal and false otherwise.
func IsTerminal(f *os.File) bool {
	return terminal.IsTerminal(int(f.Fd()))
}
