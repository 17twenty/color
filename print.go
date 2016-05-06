package color

import (
	"fmt"
	"io"
	"os"
)

var stdout = NewPrinter(os.Stdout, IsTerminal(os.Stdout))

// Printf formats according to a format specifier or highlight verb and writes to standard
// output. It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return stdout.Printf(format, a...)
}

// Printfp is the same as Printf but takes a prepared format struct.
func Printfp(f *Format, a ...interface{}) (n int, err error) {
	return stdout.Printfp(f, a...)
}

// Printer prints to a writer using highlight verbs.
type Printer struct {
	w     io.Writer // underlying writer
	color bool      // dictates whether highlight verbs are processed or stripped
}

// NewPrinter creates a new Printer that writes to out.
// The color argument dictates whether color output is enabled.
func NewPrinter(out io.Writer, color bool) *Printer {
	return &Printer{out, color}
}

// Printf first processes the highlight verbs in format and then calls
// fmt.Fprintf with the processed format and the other arguments.
// It returns the number of bytes written an any write error encountered.
func (p *Printer) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, Run(format, p.color), a...)
}

// Printfp is the same as p.Printf but takes a prepared format struct.
func (p *Printer) Printfp(f *Format, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, f.Get(p.color), a...)
}
