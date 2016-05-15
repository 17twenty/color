/*
Package log implements a simple logging package with support for highlight verbs.
It defines a Logger type with methods for formatting and printing output.

It also defines a global standard Logger that writes to standard error. Color output
will only be enabled if standard error is a terminal.
Use the helper functions Printf, Fatalf, Panicf, SetOutput and SetColor to access it.
*/
package log

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/nhooyr/color"
)

// Logger is a very simple logger similar to log.Logger but it supports the highlight verbs.
type Logger struct {
	mu    sync.Mutex // ensures atomic writes
	out   io.Writer  // destination for output
	color bool       // enable color output
}

// New creates a new Logger. The out argument sets the
// destination to which log data will be written.
// The color argument dictates whether color output is enabled.
func New(out io.Writer, color bool) *Logger {
	return &Logger{out: out, color: color}
}

// Printf expands f to its appropiate string and then calls fmt.Fprintf with the resulting
// string and the variadic arguments to write to out.
// It will expand each Format in a to its appropiate string before calling fmt.Fprintf.
func (l *Logger) Printf(f interface{}, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(a, l.color)
	switch v := f.(type) {
	case string:
		fmt.Fprintf(l.out, color.Run(v, l.color), a...)
	case *color.Format:
		fmt.Fprintf(l.out, v.Get(l.color), a...)
	default:
		panic(color.ErrBadFormat)
	}
}

// Print calls fmt.Fprint to print to the underlying writer.
// It will expand each Format in a to its appropiate string before calling fmt.Fprint.
func (l *Logger) Print(a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(a, l.color)
	fmt.Fprint(l.out, a...)
}

// Println calls fmt.Fprintln to print to the underlying writer.
// It will expand each Format in a to its appropiate string before calling fmt.Fprintln.
func (l *Logger) Println(a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(a, l.color)
	fmt.Fprintln(l.out, a...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(f interface{}, a ...interface{}) {
	l.Printf(f, a...)
	os.Exit(1)
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(a ...interface{}) {
	l.mu.Lock()
	color.ExpandFormats(a, l.color)
	fmt.Fprint(l.out, a...)
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(a ...interface{}) {
	l.mu.Lock()
	color.ExpandFormats(a, l.color)
	fmt.Fprintln(l.out, a...)
	os.Exit(1)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(f interface{}, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(a, l.color)
	var s string
	switch v := f.(type) {
	case string:
		s = fmt.Sprintf(color.Run(v, l.color), a...)
	case *color.Format:
		s = fmt.Sprintf(v.Get(l.color), a...)
	default:
		panic(color.ErrBadFormat)
	}
	io.WriteString(l.out, s)
	panic(s)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(a, l.color)
	s := fmt.Sprint(a...)
	io.WriteString(l.out, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(an ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(an, l.color)
	s := fmt.Sprintln(an...)
	io.WriteString(l.out, s)
	panic(s)
}

// SetOutput sets the output destination.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

// SetColor sets whether colored output is enabled.
func (l *Logger) SetColor(color bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.color = color
}

var std = New(os.Stderr, color.IsTerminal(os.Stderr))

// Printf calls the standard Logger's Printf method.
func Printf(f interface{}, a ...interface{}) {
	std.Printf(f, a...)
}

// Print calls the standard Logger's Printf method.
func Print(a ...interface{}) {
	std.Print(a...)
}

// Println calls the standard Logger's Println method.
func Println(a ...interface{}) {
	std.Println(a...)
}

// Fatalf calls the standard Logger's Fatalf method.
func Fatalf(f interface{}, a ...interface{}) {
	std.Fatalf(f, a...)
}

// Fatal calls the standard Logger's Fatal method.
func Fatal(a ...interface{}) {
	std.Fatal(a...)
}

// Fatalln calls the standard Logger's Fatalln method.
func Fatalln(a ...interface{}) {
	std.Fatalln(a...)
}

// Panicf calls the standard Logger's Panicf method.
func Panicf(f interface{}, a ...interface{}) {
	std.Panicf(f, a...)
}

// Panic calls the standard Logger's Panic method.
func Panic(a ...interface{}) {
	std.Panic(a...)
}

// Panicln calls the standard Logger's Panicln method.
func Panicln(a ...interface{}) {
	std.Panicln(a...)
}

// SetOutput sets the output destination of the standard Logger.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

// SetColor sets whether colored output is enabled for the standard Logger.
func SetColor(color bool) {
	std.SetColor(color)
}
