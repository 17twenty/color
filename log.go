package color

import (
	"io"
	"log"
)

type ColorLogger struct {
	*log.Logger
}

func (l *ColorLogger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(Scolorf(format), v...)
}

func (l *ColorLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(Scolorf(format), v...)
}

func (l *ColorLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(Scolorf(format), v...)
}

func NewLogger(out io.Writer, prefix string, flag int) *ColorLogger {
	return &ColorLogger{log.New(out, prefix, flag)}
}
