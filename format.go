package color

import (
	"fmt"
	"strings"
)

// Format represents a format string with the highlight verbs fully parsed.
type Format struct {
	colored  string // highlight verbs replaced with their escape sequences
	stripped string // highlight verbs stripped
}

// Get returns the colored string if color is true, and the stripped string otherwise.
func (f *Format) Get(color bool) string {
	if color {
		return f.colored
	}
	return f.stripped
}

// Append appends f2's strings to f's and then returns the resulting Format.
func (f *Format) Append(f2 *Format) *Format {
	return &Format{f.colored + f2.colored, f.stripped + f2.stripped}
}

// AppendString appends s to f's strings and then returns the resulting Format.
func (f *Format) AppendString(s string) *Format {
	return &Format{f.colored + s, f.stripped + s}
}

// Eprintf calls fmt.Sprintf using f's strings and the rest of the arguments.
// It then returns the resulting Format.
func (f *Format) Eprintf(a ...interface{}) *Format {
	return &Format{fmt.Sprintf(f.colored, a...), fmt.Sprintf(f.stripped, a...)}
}

// Insert replaces "%a" in f's strings with f2's strings.
func (f *Format) Insert(f2 *Format) *Format {
	return &Format{strings.Replace(f.colored, "%a", f2.colored, 1), strings.Replace(f.stripped, "%a", f2.stripped, 1)}
}

// InsertEmpty replaces "%a" in f's strings with "">
func (f *Format) InsertEmpty() *Format {
	return &Format{strings.Replace(f.colored, "%a", "", 1), strings.Replace(f.stripped, "%a", "", 1)}
}

// Prepare returns a Format structure using f as the base string.
func Prepare(f string) *Format {
	return &Format{Highlight(f), Strip(f)}
}
