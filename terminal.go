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

// IsTerminalWriter returns true if w is a terminal and false otherwise.
func IsTerminalWriter(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return terminal.IsTerminal(int(f.Fd()))
	}
	return false
}
