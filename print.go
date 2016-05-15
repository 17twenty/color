package color

import (
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

var ErrBadFormat = errors.New("color: format is not a string or *color.Format ")

// Printer prints to a writer using highlight verbs.
type Printer struct {
	out   io.Writer // underlying writer
	color bool      // dictates whether highlight verbs are processed or stripped
}

// New creates a new Printer that writes to out.
// The color argument dictates whether color output is enabled.
func New(out io.Writer, color bool) *Printer {
	return &Printer{out, color}
}

// Printf expands f to its appropiate string and then calls fmt.Fprintf with the resulting
// string and the variadic arguments to write to out.
// It will expand each Format in v to its appropiate string before calling fmt.Fprintf.
// It returns the number of bytes written an any write error encountered.
func (p *Printer) Printf(f interface{}, a ...interface{}) (n int, err error) {
	ExpandFormats(a, p.color)
	switch v := f.(type) {
	case string:
		return fmt.Fprintf(p.out, Run(v, p.color), a...)
	case *Format:
		return fmt.Fprintf(p.out, v.Get(p.color), a...)
	}
	panic(ErrBadFormat)
}

// Print calls fmt.Fprint to print to the underlying writer.
// It will expand each Format in a to its appropiate string before calling Print.
func (p *Printer) Print(a ...interface{}) (n int, err error) {
	ExpandFormats(a, p.color)
	return fmt.Fprint(p.out, a...)
}

// Println calls fmt.Fprintln to print to the underlying writer.
// It will expand each Format in a to its appropiate string before calling Println.
func (p *Printer) Println(a ...interface{}) (n int, err error) {
	ExpandFormats(a, p.color)
	return fmt.Fprintln(p.out, a...)
}

var std = New(os.Stdout, IsTerminal(os.Stdout))

// Printf calls the standard output Printer's Printf method.
func Printf(f *Format, a ...interface{}) (n int, err error) {
	return std.Printf(f, a...)
}

// Print calls the standard output Printer's Print method.
func Print(a ...interface{}) (n int, err error) {
	return std.Print(a...)
}

// Println calls the standard output Printer's Println method.
func Println(a ...interface{}) (n int, err error) {
	return std.Println(a...)
}

// IsTerminal returns true if f is a terminal and false otherwise.
func IsTerminal(f *os.File) bool {
	return terminal.IsTerminal(int(f.Fd()))
}
