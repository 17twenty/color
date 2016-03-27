package color

import "fmt"

func Sprintf(format string, a ...interface{}) string {
	h := &highlighter{s: format}
	h.run()
	return fmt.Sprintf(h.s, a...)
}
