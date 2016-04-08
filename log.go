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

func (l *Logger) highlight(s string) string {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.color {
		return Highlight(s)
	}
	return stripVerbs(s)
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

// SetOutput first checks if w is a terminal, then sets the output.
func (l *Logger) SetOutput(w io.Writer) {
	l.isTerminal(w)
	l.Logger.SetOutput(w)
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
	l.Logger.SetPrefix(l.highlight(prefix))
}

// SetColor sets whether the Logger should output in color.
func (l *Logger) SetColor(b bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	// fix underlying logger prefix.
	if b {
		l.Logger.SetPrefix(Highlight(l.prefix))
	} else {
		l.Logger.SetPrefix(stripVerbs(l.prefix))
	}
	l.color = b
}

// isTerminal turns on color output if w is a terminal.
func (l *Logger) isTerminal(w io.Writer) {
	if isTerminal(w) {
		l.SetColor(true)
	} else {
		l.SetColor(false)
	}
}

// NewLogger creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line
// and it can contain highlighting verbs.
// The flag argument defines the logging properties.
func NewLogger(out io.Writer, prefix string, flag int) (l *Logger) {
	l = &Logger{Logger: log.New(out, "", flag), prefix: prefix}
	l.isTerminal(out)
	return
}
