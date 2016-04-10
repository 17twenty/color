package color

type Format struct {
	Highlighted string
	Stripped    string
}

func Prepare(format string) (f *Format) {
	return &Format{Shighlightf(format), Sstripf(format)}
}

func (f *Format) Get(color bool) string {
	if color {
		return f.Highlighted
	}
	return f.Stripped
}
