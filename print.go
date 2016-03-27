package color

import (
	"fmt"
	"io"
	"os"
)

func scolorf(s string) string {
	h := &highlighter{s: s}
	h.run()
	return h.s
}

func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, scolorf(format), a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(scolorf(format), a...)
}
