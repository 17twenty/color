package color

import (
	"fmt"
	"io"
	"os"
)

func Sprint(format string) string {
	h := &highlighter{s: format}
	h.run()
	return h.s
}

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, Sprint(format), a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(Sprint(format), a...)
}
