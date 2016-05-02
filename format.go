package color

import (
	"fmt"
)

type Format struct {
	colored  string
	stripped string
}

func (f *Format) Append(f2 *Format) *Format {
	return &Format{f.colored + f2.colored, f.stripped + f2.stripped}
}

func (f *Format) Get(color bool) string {
	if color {
		return f.colored
	}
	return f.stripped
}

// TODO: add a method like this that takes a Format and Printf's it into a new Format using the
// current format's strings as the format string.
func (f *Format) Printf(a ...interface{}) *Format {
	return &Format{fmt.Sprintf(f.colored, a...), fmt.Sprintf(f.stripped, a...)}
}

func Prepare(f string) *Format {
	return &Format{Highlight(f), Strip(f)}
}
