// Bridge package to expose color internals to tests in the color_test
// package.

package color

var (
	Colors = colors
	Ti     = ti
	Modes = modes
)

const (
	ErrInvalid = errInvalid
	ErrMissing = errMissing
	ErrBadAttr = errBadAttr
	ErrShort = errShort
)
