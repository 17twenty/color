package color

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// Printer prints to a writer using highlight verbs.
type Printer struct {
	out   io.Writer // underlying writer
	color bool      // enable color output
}

// New creates a new Printer that writes to out.
// The color argument dictates whether color output is enabled.
func New(out io.Writer, color bool) *Printer {
	return &Printer{out, color}
}

// Printf first processes the highlight verbs in format and then calls
// fmt.Fprintf with the processed format and the other arguments.
// It will expand each Format in a to its appropriate string before calling fmt.Fprintf.
// It returns the number of bytes written an any write error encountered.
func (p *Printer) Printf(format string, a ...interface{}) (n int, err error) {
	ExpandFormats(p.color, a)
	return fmt.Fprintf(p.out, Run(format, p.color), a...)
}

// Printfp is the same as p.Printf but takes a prepared format struct.
func (p *Printer) Printfp(f *Format, a ...interface{}) (n int, err error) {
	ExpandFormats(p.color, a)
	return fmt.Fprintf(p.out, f.Get(p.color), a...)
}

// Print calls fmt.Fprint to print to the underlying writer.
// It will expand each Format in a to its appropriate string before calling fmt.Fprint.
func (p *Printer) Print(a ...interface{}) (n int, err error) {
	ExpandFormats(p.color, a)
	return fmt.Fprint(p.out, a...)
}

// Println calls fmt.Fprintln to print to the underlying writer.
// It will expand each Format in a to its appropriate string before calling fmt.Fprintln.
func (p *Printer) Println(a ...interface{}) (n int, err error) {
	ExpandFormats(p.color, a)
	return fmt.Fprintln(p.out, a...)
}

// IsTerminal returns true if f is a terminal and false otherwise.
func IsTerminal(f *os.File) bool {
	return terminal.IsTerminal(int(f.Fd()))
}

var std = New(os.Stdout, IsTerminal(os.Stdout))

// Printf calls the standard output Printer's Printf method.
func Printf(format string, a ...interface{}) (n int, err error) {
	return std.Printf(format, a...)
}

// Printfp calls the standard output Printer's Printfp method.
func Printfp(f *Format, a ...interface{}) (n int, err error) {
	return std.Printfp(f, a...)
}

// Print calls the standard output Printer's Print method.
func Print(a ...interface{}) (n int, err error) {
	return std.Print(a...)
}

// Println calls the standard output Printer's Println method.
func Println(a ...interface{}) (n int, err error) {
	return std.Println(a...)
}
