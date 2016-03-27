package color

import (
	"io"
	"log"
)

type ColorLogger struct {
	*log.Logger
}

func (l *ColorLogger) Printf(format string, v ...interface{}) {
	l.Logger.Printf(scolorf(format), v...)
}

func (l *ColorLogger) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatalf(scolorf(format), v...)
}

func (l *ColorLogger) Panicf(format string, v ...interface{}) {
	l.Logger.Panicf(scolorf(format), v...)
}

func NewLogger(out io.Writer, prefix string, flag int) *ColorLogger {
	return &ColorLogger{log.New(out, prefix, flag)}
}
