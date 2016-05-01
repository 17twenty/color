package color

import (
	"io"
	"log"
)

// Logger is a very thin wrapper around log.Logger to support the highlighting verbs.
type Logger struct {
	*log.Logger      // TODO unexport this somehow
	color       bool // dictates whether highlight verbs are processed or stripped
}

// NewLogger creates a new Logger. The out argument sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// It can contain highlight verbs, however l.SetPrefix does not allow highlight verbs.
// The flag argument defines the logging properties.
// The color argument dictates whether color output is enabled.
func NewLogger(out io.Writer, prefix string, flag int, color bool) *Logger {
	return &Logger{log.New(out, Run(prefix, color), flag), color}
}

// Printfh first calls l.Prepare to process the highlight verbs and then
// calls l.Logger.Printf to print to the logger.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(Run(format, l.color), v...)
}

// Fatalfh is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(Run(format, l.color), v...)
}

// Panicfh is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(Run(format, l.color), v...)
}

// Printfh first calls l.Prepare to process the highlight verbs and then
// calls l.Logger.Printf to print to the logger.
func (l *Logger) Eprintf(f *Format, v ...interface{}) {
	l.Logger.Printf(f.Get(l.color), v...)
}

// Fatalfh is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Efatalf(f *Format, v ...interface{}) {
	l.Logger.Fatalf(f.Get(l.color), v...)
}

// Panicfh is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Epanicf(f *Format, v ...interface{}) {
	l.Logger.Panicf(f.Get(l.color), v...)
}

// SetOutput panics if called. You cannot change the output writer once the Logger is created.
// TODO fix this.
func (l *Logger) SetOutput(w io.Writer) {
	panic("SetOutput is not supported on *color.Logger")
}
