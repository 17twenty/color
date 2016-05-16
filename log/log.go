/*
Package log implements a simple logging package with support for highlight verbs.
It defines a Logger type with methods for formatting and printing output.

It also defines a global standard Logger that writes to standard error. Color output
will only be enabled if standard error is a terminal.
Use the helper functions Printf[p], Fatalf[p], Panicf[p], and SetOutput to access it.
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
// It will expand each Format in v to its appropiate string before calling fmt.Fprintf.
func (l *Logger) Printf(f *color.Format, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Replace(l.color, v)
	fmt.Fprintf(l.out, f.Get(l.color), v...)
}

// Print calls fmt.Fprint to print to the underlying writer.
// It will expand each Format in v to its appropiate string before calling fmt.Fprint.
func (l *Logger) Print(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Replace(l.color, v)
	fmt.Fprint(l.out, v...)
}

// Println calls fmt.Fprintln to print to the underlying writer.
// It will expand each Format in v to its appropiate string before calling fmt.Fprintln.
func (l *Logger) Println(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Replace(l.color, v)
	fmt.Fprintln(l.out, v...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(f *color.Format, v ...interface{}) {
	l.mu.Lock()
	color.Replace(l.color, v)
	fmt.Fprintf(l.out, f.Get(l.color), v...)
	os.Exit(1)
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.mu.Lock()
	color.Replace(l.color, v)
	fmt.Fprint(l.out, v...)
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	l.mu.Lock()
	color.Replace(l.color, v)
	fmt.Fprintln(l.out, v...)
	os.Exit(1)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(f *color.Format, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Replace(l.color, v)
	s := fmt.Sprintf(f.Get(l.color), v...)
	io.WriteString(l.out, s)
	panic(s)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Replace(l.color, v)
	s := fmt.Sprint(v...)
	io.WriteString(l.out, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.Replace(l.color, v)
	s := fmt.Sprintln(v...)
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
func Printf(f *color.Format, v ...interface{}) {
	std.Printf(f, v...)
}

// Print calls the standard Logger's Printf method.
func Print(v ...interface{}) {
	std.Print(v...)
}

// Println calls the standard Logger's Println method.
func Println(v ...interface{}) {
	std.Println(v...)
}

// Fatalf calls the standard Logger's Fatalf method.
func Fatalf(f *color.Format, v ...interface{}) {
	std.Fatalf(f, v...)
}

// Fatal calls the standard Logger's Fatal method.
func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

// Fatalln calls the standard Logger's Fatalln method.
func Fatalln(v ...interface{}) {
	std.Fatalln(v...)
}

// Panicf calls the standard Logger's Panicf method.
func Panicf(f *color.Format, v ...interface{}) {
	std.Panicf(f, v...)
}

// Panic calls the standard Logger's Panic method.
func Panic(v ...interface{}) {
	std.Panic(v...)
}

// Panicln calls the standard Logger's Panicln method.
func Panicln(v ...interface{}) {
	std.Panicln(v...)
}

// SetOutput sets the output destination of the standard Logger.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

// SetColor sets whether colored output is enabled for the standard Logger.
func SetColor(color bool) {
	std.SetColor(color)
}
