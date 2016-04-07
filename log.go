package color

import (
	"io"
	"log"
)

// Logger wraps a log.Logger's Printf functions to support the highlighting verbs.
type Logger struct {
	*log.Logger
}

// Printf calls l.Logger.Printf to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(Highlight(format), v...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(Highlight(format), v...)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(Highlight(format), v...)
}

// NewLogger creates a new ColorLogger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
func NewLogger(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{log.New(out, prefix, flag)}
}
