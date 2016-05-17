/*
Package log implements a simple logging package with support for highlight verbs.
It defines a Logger type with methods for formatting and printing output.

It also defines a global standard Logger that writes to standard error. Color output
will only be enabled if standard error is a terminal.
Use the helper functions Print[f|ln|p], Fatal[f|ln|p], Panicf[f|ln|p], SetOutput and SetColor to access it.
*/
package log

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/nhooyr/color"
)

// Logger is a very simple logger, similar to log.logger but it supports highlight verbs.
type Logger struct {
	mu    sync.Mutex  // ensures atomic writes
	out   *lineWriter // ensures output is written on separate lines
	color bool        // enable color output
}

// New creates a new Logger. The out argument sets the
// destination to which log data will be written.
// The color argument dictates whether color output is enabled.
func New(out io.Writer, color bool) *Logger {
	return &Logger{out: &lineWriter{out}, color: color}
}

// Printf processes the highlight verbs in format and then calls
// fmt.Fprintf to print to the underlying writer.
// It will expand each Format in v to its appropiate string before calling fmt.Fprintf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(l.color, v)
	fmt.Fprintf(l.out, color.Run(format, l.color), v...)
}

// Printfp is the same as l.Printf but takes a prepared format struct.
func (l *Logger) Printfp(f *color.Format, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(l.color, v)
	fmt.Fprintf(l.out, f.Get(l.color), v...)
}

// Print calls fmt.Fprint to print to the underlying writer.
// It will expand each Format in v to its appropiate string before calling fmt.Fprint.
func (l *Logger) Print(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(l.color, v)
	fmt.Fprint(l.out, v...)
}

// Println calls fmt.Fprintln to print to the underlying writer.
// It will expand each Format in v to its appropiate string before calling fmt.Fprintln.
func (l *Logger) Println(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(l.color, v)
	fmt.Fprintln(l.out, v...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.mu.Lock()
	color.ExpandFormats(l.color, v)
	fmt.Fprintf(l.out, color.Run(format, l.color), v...)
	os.Exit(1)
}

// Fatalfp is the same as l.Fatalf but takes a prepared format struct.
func (l *Logger) Fatalfp(f *color.Format, v ...interface{}) {
	l.mu.Lock()
	color.ExpandFormats(l.color, v)
	fmt.Fprintf(l.out, f.Get(l.color), v...)
	os.Exit(1)
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Logger) Fatal(v ...interface{}) {
	l.mu.Lock()
	color.ExpandFormats(l.color, v)
	fmt.Fprint(l.out, v...)
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Logger) Fatalln(v ...interface{}) {
	l.mu.Lock()
	color.ExpandFormats(l.color, v)
	fmt.Fprintln(l.out, v...)
	os.Exit(1)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(l.color, v)
	s := fmt.Sprintf(format, v...)
	l.out.WriteString(s)
	panic(s)
}

// Panicfp is the same as l.Panicf but takes a prepared format struct.
func (l *Logger) Panicfp(f *color.Format, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(l.color, v)
	s := fmt.Sprintf(f.Get(l.color), v...)
	l.out.WriteString(s)
	panic(s)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(l.color, v)
	s := fmt.Sprint(v...)
	l.out.WriteString(s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Logger) Panicln(v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	color.ExpandFormats(l.color, v)
	s := fmt.Sprintln(v...)
	l.out.WriteString(s)
	panic(s)
}

// SetOutput sets the output destination.
func (l *Logger) SetOutput(out io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out.w = out
}

// SetColor sets whether colored output is enabled.
func (l *Logger) SetColor(color bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.color = color
}

// lineWriter ensures that each Write to the underlying writer will end on a newline.
type lineWriter struct {
	w io.Writer // underlying writer
}

// Write writes to the underlying writer but ensures that the write ends on a newline.
func (lw *lineWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 || p[len(p)-1] != '\n' {
		return lw.w.Write(append(p, '\n'))
	}
	return lw.w.Write(p)
}

// WriteString is the same as lw.Write but takes a string.
func (lw *lineWriter) WriteString(s string) (n int, err error) {
	if len(s) == 0 || s[len(s)-1] == '\n' {
		p := make([]byte, len(s))
		copy(p, s)
		p[len(s)] = '\n'
		return lw.w.Write(p)
	}
	return io.WriteString(lw.w, s)
}

var std = New(os.Stderr, color.IsTerminal(os.Stderr))

// Printf calls the standard Logger's Printf method.
func Printf(format string, v ...interface{}) {
	std.Printf(format, v...)
}

// Printfp calls the standard Logger's Printfp method.
func Printfp(f *color.Format, v ...interface{}) {
	std.Printfp(f, v...)
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
func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

// Fatalfp calls the standard Logger's Fatalfp method.
func Fatalfp(f *color.Format, v ...interface{}) {
	std.Fatalfp(f, v...)
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
func Panicf(format string, v ...interface{}) {
	std.Panicf(format, v...)
}

// Panicfp calls the standard Logger's Panicfp method.
func Panicfp(f *color.Format, v ...interface{}) {
	std.Panicfp(f, v...)
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
