package color

import (
	"fmt"
	"io"
	"os"
)

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// It only prints in color if the writer is a terminal, otherwise it prints normally.
// Use a color.Printer if you want full control over when to print in color or you want
// to avoid the repetitive terminal checks.
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if IsTerminal(w) {
		return fmt.Fprintf(w, Shighlightf(format), a...)
	}
	return fmt.Fprintf(w, Sstripf(format), a...)
}

// Efprintf is the same as Fprintf but takes a prepared Format object.
func Efprintf(w io.Writer, f *Format, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, f.Get(IsTerminal(w)), a...)
}

var stdout = NewPrinter(os.Stdout, PerformCheck)

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return stdout.Printf(format, a...)
}

// Eprintf is the same as Printf but takes a prepared Format object.
func Eprintf(f *Format, a ...interface{}) (n int, err error) {
	return stdout.Eprintf(f, a...)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(Shighlightf(format), a...)
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
	PerformCheck = iota // check if a terminal, and if so enable colored output
	EnableColor         // enable colored output
	DisableColor        // disable colored output
)

// NewPrinter creates a new Printer. It enables colored output
// based on the flag.
func NewPrinter(out io.Writer, flag uint8) (p *Printer) {
	p = &Printer{w: out}
	if flag == PerformCheck && IsTerminal(out) || flag == EnableColor {
		p.color = true
	}
	return
}

// Printf calls fmt.Fprintf to print to the writer.
// Arguments are handled in the manner of color.Printf.
func (p *Printer) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, Scolorf(format, p.color), a...)
}

// Eprintf is the same as p.Printf but takes a prepared Format object.
func (p *Printer) Eprintf(f *Format, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, f.Get(p.color), a...)
}
