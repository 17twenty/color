package color

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var std = New(os.Stdout, IsTerminal(os.Stdout))

// Printf calls the standard Printer's Printf method.
func Printf(format string, a ...interface{}) (n int, err error) {
	return std.Printf(format, a...)
}

// Printfp calls the standard Printer's Printfp method.
func Printfp(f *Format, a ...interface{}) (n int, err error) {
	return std.Printfp(f, a...)
}

// Print calls the standard Printer's Print method.
func Print(a ...interface{}) (n int, err error) {
	return std.Print(a...)
}

// Println calls the standard Printer's Println method.
func Println(a ...interface{}) (n int, err error) {
	return std.Println(a...)
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
// It will expand each Format in a to its appropiate string before calling Printf.
// It returns the number of bytes written an any write error encountered.
func (p *Printer) Printf(format string, a ...interface{}) (n int, err error) {
	Replace(p.color, a)
	return fmt.Fprintf(p.w, Run(format, p.color), a...)
}

// Printfp is the same as p.Printf but takes a prepared format struct.
func (p *Printer) Printfp(f *Format, a ...interface{}) (n int, err error) {
	Replace(p.color, a)
	return fmt.Fprintf(p.w, f.Get(p.color), a...)
}

// Print calls fmt.Fprint to print to the underlying writer.
// It will expand each Format in a to its appropiate string before calling Print.
func (p *Printer) Print(a ...interface{}) (n int, err error) {
	Replace(p.color, a)
	return fmt.Fprint(p.w, a...)
}

// Println calls fmt.Fprintln to print to the underlying writer.
// It will expand each Format in a to its appropiate string before calling Println.
func (p *Printer) Println(a ...interface{}) (n int, err error) {
	Replace(p.color, a)
	return fmt.Fprintln(p.w, a...)
}

// IsTerminal returns true if f is a terminal and false otherwise.
func IsTerminal(f *os.File) bool {
	return terminal.IsTerminal(int(f.Fd()))
}
