package color

import "fmt"

// Format represents a format string with the highlight verbs fully parsed.
type Format struct {
	colored  string // highlight verbs replaced with their escape sequences
	stripped string // highlight verbs stripped
}

// Prepare returns a Format structure using f as the base string.
func Prepare(f string) *Format {
	return &Format{Highlight(f), Strip(f)}
}

// Get returns the colored string if color is true, and the stripped string otherwise.
func (f *Format) Get(color bool) string {
	if color {
		return f.colored
	}
	return f.stripped
}

// Eprintf calls fmt.Sprintf using f's strings and the rest of the arguments.
// It will expand each Format in a to its appropiate string before calling Sprintf.
// It then returns the resulting Format.
func (f *Format) Eprintf(a ...interface{}) *Format {
	m := make(map[int]*Format)
	for i, v := range a {
		if f, ok := v.(*Format); ok {
			a[i] = f.Get(true)
			m[i] = f
		}
	}
	rf := new(Format)
	rf.colored = fmt.Sprintf(f.colored, a...)
	for i, f := range m {
		a[i] = f.Get(false)
	}
	rf.stripped = fmt.Sprintf(f.stripped, a...)
	return rf
}

// Replace replaces each Format in a with its appropiate string according to color.
func Replace(color bool, a []interface{}) {
	for i, v := range a {
		if f, ok := v.(*Format); ok {
			a[i] = f.Get(color)
		}
	}
}
