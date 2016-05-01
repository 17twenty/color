package color

import (
	"fmt"
	"io"
	"os"
)

// Fprintfh formats according to a format specifier or highlight verb and writes to w.
// It returns the number of bytes written and any write error encountered. It only prints in
// color if the writer is a terminal, otherwise it prints normally. Use a Printer if you want
// full control over when to print in color or you want to avoid the repetitive terminal checks.
func Fprintfh(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, Run(format, IsTerminal(w)))
}

var stdout = NewPrinter(os.Stdout, IsTerminal(os.Stdout))

// Printfh formats according to a format specifier or highlight verb and writes to standard
// output. It returns the number of bytes written and any write error encountered.
func Printfh(format string, a ...interface{}) (n int, err error) {
	return stdout.Printfh(format, a...)
}

// Prepare returns the format string with only the highlight verbs processed.
func Prepare(format string) string {
	return stdout.Prepare(format)
}

// Printer prints to a writer. Use this over Fprintf when you want full control over when
// to color the output or you want to avoid the repetitive terminal checks done by Fprinf.
type Printer struct {
	w     io.Writer
	color bool // dictates whether highlight verbs are processed or stripped
}

// NewPrinter creates a new Printer.
// The color argument dictates whether color output is enabled.
func NewPrinter(out io.Writer, color bool) *Printer {
	return &Printer{out, color}
}

// Printf calls fmt.Fprintf to print to the writer.
func (p *Printer) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, format, a...)
}

// Printfh first calls p.Prepare to process the highlight verbs and then
// calls fmt.Fprintf to print to the writer.
func (p *Printer) Printfh(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, p.Prepare(format), a...)
}

// Prepare returns the format string with the highlight verbs processed.
// It is a thin wrapper around Run.
func (p *Printer) Prepare(format string) string {
	return Run(format, p.color)
}
