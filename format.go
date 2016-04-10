package color

// Format represents a processed color string.
type Format struct {
	highlighted string // colored string
	stripped    string // highlight verbs stripped.
}

// Prepare takes a format string and returns a Format object representing
// the processed format string.
func Prepare(format string) (f *Format) {
	return &Format{Shighlightf(format), Sstripf(format)}
}

// Get returns the colored version of the processed string if color is true.
// Otherwise it returns the stripped version.
func (f *Format) Get(color bool) string {
	if color {
		return f.highlighted
	}
	return f.stripped
}
