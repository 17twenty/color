package color

// scolorf returns the highlighted string if color is true, otherwise it returns the string with the highlight verbs stripped.
func scolorf(format string, color bool) string {
	if color {
		return Shighlightf(format)
	}
	return Sstripf(format)
}
