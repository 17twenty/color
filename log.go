package color

import (
	"io"
	"log"
	"sync"
)

// Logger wraps a log.Logger's Printf functions to support the highlighting verbs.
type Logger struct {
	*log.Logger
	prefix string
	color  bool // dictates if highlight verbs are applied
	mu     sync.Mutex
}

// NewLogger creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line
// and it can contain highlighting verbs.
// The flag argument defines the logging properties.
// It checks if the writer is a terminal and enables color output accordingly.
func NewLogger(out io.Writer, prefix string, flag int) (l *Logger) {
	l = &Logger{Logger: log.New(out, "", flag), prefix: prefix}
	l.isTerminal(out)
	return
}

// scolorf is a convenience function for highlighting strings according to
// whether color output is set.
func (l *Logger) scolorf(s string) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return scolorf(s, l.color)
}

// Printf calls l.Logger.Printf to print to the logger.
// Arguments are handled in the manner of color.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(l.scolorf(format), v...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(l.scolorf(format), v...)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(l.scolorf(format), v...)
}

// Eprintf is the same as l.Panicf but takes a prepared Format object.
func (l *Logger) Eprintf(f *Format, v ...interface{}) {
	l.Logger.Printf(f.Get(l.color), v...)
}

// Efatalf is the same as l.Fatalf but takes a prepared Format object.
func (l *Logger) Efatalf(f *Format, v ...interface{}) {
	l.Logger.Fatalf(f.Get(l.color), v...)
}

// Epanicf is the same as l.Panicf but takes a prepared Format object.
func (l *Logger) Epanicf(f *Format, v ...interface{}) {
	l.Logger.Panicf(f.Get(l.color), v...)
}

// SetOutput checks if the writer is a terminal and sets the color output accordingly.
// Then it sets the underlying loggers output to the writer.
func (l *Logger) SetOutput(w io.Writer) {
	l.isTerminal(w)
	l.Logger.SetOutput(w)
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
	l.Logger.SetPrefix(l.scolorf(prefix))
}

// EnableColor enables color output.
func (l *Logger) EnableColor() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.color = true
	l.Logger.SetPrefix(Shighlightf(l.prefix))
}

// DisableColor disables color output.
func (l *Logger) DisableColor() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.color = false
	l.Logger.SetPrefix(Sstripf(l.prefix))
}

// isTerminal turns on color output if w is a terminal.
func (l *Logger) isTerminal(w io.Writer) {
	if IsTerminal(w) {
		l.EnableColor()
	} else {
		l.DisableColor()
	}
}
