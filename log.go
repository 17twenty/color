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

// highlight is a convenience function for highlighting strings according to
// whether color output is set.
func (l *Logger) highlight(s string) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.color {
		return shighlightf(s)
	}
	return sstripf(s)
}

// Printf calls l.Logger.Printf to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(l.highlight(format), v...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(l.highlight(format), v...)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(l.highlight(format), v...)
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
	l.Logger.SetPrefix(l.highlight(prefix))
}

// EnableColor enables color output.
func (l *Logger) EnableColor() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.color = true
	l.Logger.SetPrefix(shighlightf(l.prefix))
}

// DisableColor disables color output.
func (l *Logger) DisableColor() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.color = false
	l.Logger.SetPrefix(sstripf(l.prefix))
}

// isTerminal turns on color output if w is a terminal.
func (l *Logger) isTerminal(w io.Writer) {
	if isTerminal(w) {
		l.EnableColor()
	} else {
		l.DisableColor()
	}
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
