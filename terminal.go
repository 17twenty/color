package color

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// IsTerminal returns true if the writer passed is a terminal and false otherwise.
func IsTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok && terminal.IsTerminal(int(f.Fd())) {
		return true
	}
	return false
}
