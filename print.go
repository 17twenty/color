package color

import (
	"fmt"
	"io"
	"os"
)

// Fprintfh formats according to a format specifier or a highlight verb and writes to w.
// It returns the number of bytes written and any write error encountered.
// It only prints in color if the writer is a terminal, otherwise it prints normally.
// Use a color.Printer if you want full control over when to print in color or you want
// to avoid the repetitive terminal checks.
func Fprintfh(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if IsTerminal(w) {
		return fmt.Fprintf(w, Highlight(format), a...)
	}
	return fmt.Fprintf(w, Strip(format), a...)
}

var stdout = NewPrinter(os.Stdout, PerformCheck)

// Printfh formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printfh(format string, a ...interface{}) (n int, err error) {
	return stdout.Printfh(format, a...)
}

// Prepare returns the format string with only the highlight verbs processed.
func Prepare(format string) string {
	return stdout.Prepare(format)
}

// Printer prints to a writer. It is exactly like color.Fprintf except use this when
// you want full control over when to color the output or you want to avoid the repetitive
// terminal checks done by color.Fprinf.
type Printer struct {
	w     io.Writer
	color bool
}

// Flags for setting colored output when creating a Printer.
const (
	PerformCheck = 1 << iota // check if a terminal, and if so enable colored output
	EnableColor              // enable colored output
	DisableColor             // disable colored output
)

// NewPrinter creates a new Printer. Color output is enabled or disabled based on the flag.
func NewPrinter(out io.Writer, flag int) (p *Printer) {
	p = &Printer{w: out}
	if flag == PerformCheck && IsTerminal(out) || flag == EnableColor {
		p.color = true
	}
	return
}

// Printf calls fmt.Fprintf to print to the writer.
// Arguments are handled in the manner of fmt.Printf.
func (p *Printer) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, format, a...)
}

// Printfh calls fmt.Fprintf to print to the writer.
// Arguments are handled in the manner of color.Printf.
func (p *Printer) Printfh(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, Run(format, p.color), a...)
}

// Prepare returns the format string with only the highlight verbs processed.
func (p *Printer) Prepare(format string) string {
	return Run(format, p.color)
}
