package color

import (
	"io"
	"log"
)

// Logger wraps a log.Logger's Printf functions to support the highlighting verbs.
type Logger struct {
	*log.Logger
	color bool // dictates if highlight verbs are applied
}

// NewLogger creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line
// and it can contain highlighting verbs.
// The flag argument defines the logging properties.
// It checks if the writer is a terminal and enables color output accordingly.
func NewLogger(out io.Writer, prefix string, flag int, cflag int) (l *Logger) {
	if cflag == PerformCheck && IsTerminal(out) || cflag == EnableColor {
		l = &Logger{Logger: log.New(out, Shighlightf(prefix), flag)}
		l.color = true
	} else {
		l = &Logger{Logger: log.New(out, Sstripf(prefix), flag)}
	}
	return
}

// Printf calls l.Logger.Printf to print to the logger.
// Arguments are handled in the manner of color.Printf.
func (l *Logger) Hprintf(format string, v ...interface{}) {
	l.Logger.Printf(Run(format, l.color), v...)
}

// Prepare returns the format string with only the highlight verbs handled.
func (l *Logger) Prepare(format string) string {
	return Run(format, l.color)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(Run(format, l.color), v...)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(Run(format, l.color), v...)
}

// SetOutput panics if called. It should never be called, color.Logger.SetOutput is not supported.
func (l *Logger) SetOutput(w io.Writer) {
	panic("color.Logger.SetOutput is not supported")
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.Logger.SetPrefix(Run(prefix, l.color))
}
