// Bridge package to expose color internals to tests in the color_test
// package.

package color

var Colors = colors

const (
	ErrMissing = errMissing
	ErrInvalid = errInvalid
	ErrBadAttr = errBadAttr
)
