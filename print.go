package color

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var std = New(os.Stdout, IsTerminal(os.Stdout))

// Printf formats according to a format specifier or highlight verb and writes to standard
// output. It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return std.Printf(format, a...)
}

// Printfp is the same as Printf but takes a prepared format struct.
func Printfp(f *Format, a ...interface{}) (n int, err error) {
	return std.Printfp(f, a...)
}

// Printer prints to a writer using highlight verbs.
type Printer struct {
	w     io.Writer // underlying writer
	color bool      // dictates whether highlight verbs are processed or stripped
}

// New creates a new Printer that writes to out.
// The color argument dictates whether color output is enabled.
func New(out io.Writer, color bool) *Printer {
	return &Printer{out, color}
}

// Printf first processes the highlight verbs in format and then calls
// fmt.Fprintf with the processed format and the other arguments.
// It returns the number of bytes written an any write error encountered.
func (p *Printer) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, Run(format, p.color), a...)
}

// Print calls fmt.Fprint to print to the underlying writer.
func (p *Printer) Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(p.w, a...)
}

// Println calls fmt.Fprintln to print to the underlying writer.
func (p *Printer) Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(p.w, a...)
}

// Printfp is the same as p.Printf but takes a prepared format struct.
func (p *Printer) Printfp(f *Format, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(p.w, f.Get(p.color), a...)
}

// IsTerminal returns true if f is a terminal and false otherwise.
func IsTerminal(f *os.File) bool {
	return terminal.IsTerminal(int(f.Fd()))
}
