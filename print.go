package color

import (
	"fmt"
	"io"
)

// Highlight replaces the highlight verbs in s with their appropriate
// control sequences and then returns the resulting string
func Highlight(s string) string {
	h := &highlighter{s: s}
	h.run()
	return h.s
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, Highlight(format), a...)
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Printf(Highlight(format), a...)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(Highlight(format), a...)
}
