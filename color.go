package color

// Scolorf highlights the format string accordingly to color and then returns it.
func Scolorf(format string, color bool) string {
	if color {
		return Shighlightf(format)
	}
	return Sstripf(format)
}
