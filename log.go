package color

import (
	"io"
	"log"
)

// Logger is a very lightweight wrapper around log.Logger to support the highlighting verbs.
type Logger struct {
	*log.Logger      // TODO unexport this somehow
	color       bool // dictates whether highlight verbs are processed or stripped
}

// NewLogger creates a new Logger. The out argument sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// It can contain highlight verbs, however l.SetPrefix does not allow highlight verbs.
// The flag argument defines the logging properties.
// The cflag argument dictates whether the color output is enabled.
func NewLogger(out io.Writer, prefix string, flag int, cflag int) (l *Logger) {
	l = &Logger{Logger: log.New(out, "", flag),
		color: cflag == PerformCheck && IsTerminal(out) || cflag == EnableColor}
	l.SetPrefix(l.Prepare(prefix))
	return
}

// Printfh calls l.Logger.Printf to print to the logger.
// Arguments are handled in the manner of color.Hprintf.
func (l *Logger) Printfh(format string, v ...interface{}) {
	l.Logger.Printf(l.Prepare(format), v...)
}

// Fatalfh is equivalent to l.Printfh() followed by a call to os.Exit(1).
func (l *Logger) Fatalfh(format string, v ...interface{}) {
	l.Logger.Fatalf(l.Prepare(format), v...)
}

// Panicfh is equivalent to l.Printfh() followed by a call to panic().
func (l *Logger) Panicfh(format string, v ...interface{}) {
	l.Logger.Panicf(l.Prepare(format), v...)
}

// Prepare returns the format string with only the highlight verbs processed.
func (l *Logger) Prepare(format string) string {
	return Run(format, l.color)
}

// SetOutput panics if called. You cannot change the output writer once created.
func (l *Logger) SetOutput(w io.Writer) {
	panic("color.Logger.SetOutput is not supported")
}
