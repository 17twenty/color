package color

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
// It only prints in color if the writer is a terminal, otherwise it prints normally.
// Use a color.Printer if you want full control over when to print in color or you want
// to avoid the repetitive terminal checks. Printer only checks when it is created.
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if isTerminal(w) {
		return fmt.Fprintf(w, Highlight(format), a...)
	}
	return fmt.Fprintf(w, stripVerbs(format), a...)
}

// Printf formats according to a format specifier and writes to standard output.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stdout, Highlight(format), a...)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(Highlight(format), a...)
}

// Printer prints to a writer. It is exactly like color.Fprintf except use this when
// you want full control over when to color the output or you want to avoid the repetitive
// terminal checks done by color.Fprinf. Printer only checks when it is created.
type Printer struct {
	w     io.Writer
	color bool // dictates if highlight verbs are applied
	mu    sync.Mutex
}

// highlight is a convenience function for highlighting strings according to
// whether color output is set.
func (p *Printer) highlight(s string) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.color {
		return Highlight(s)
	}
	return stripVerbs(s)
}

// Printf calls fmt.Fprintf to print to the writer.
func (p *Printer) Printf(format string, v ...interface{}) {
	fmt.Fprintf(p.w, p.highlight(format), v...)
}

// EnableColor enables color output.
func (p *Printer) EnableColor() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.color = true
}

// DisableColor disables color output.
func (p *Printer) DisableColor() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.color = false
}

// NewPrinter creates a new Printer. It checks if out is a terminal, and enables
// color output accordingly.
func NewPrinter(out io.Writer) (p *Printer) {
	p = &Printer{w: out}
	if isTerminal(out) {
		p.EnableColor()
	} else {
		p.DisableColor()
	}
	return
}
