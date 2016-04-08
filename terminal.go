package color

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok && terminal.IsTerminal(int(f.Fd())) {
		return true
	}
	return false
}
