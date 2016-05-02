package color

import (
	"fmt"
)

// TODO docs
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

func (f *Format) Printf(a ...interface{}) *Format {
	return &Format{fmt.Sprintf(f.colored, a...), fmt.Sprintf(f.stripped, a...)}
}

func (f *Format) Eprintf(f2 *Format) *Format {
	return &Format{fmt.Sprintf(f.colored, f2.colored), fmt.Sprintf(f.stripped, f2.stripped)}
}

func Prepare(f string) *Format {
	return &Format{Highlight(f), Strip(f)}
}
