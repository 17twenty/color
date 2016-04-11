package color_test

import (
	"testing"

	"github.com/nhooyr/color"
)

func TestPrepare(t *testing.T) {
	p := color.Prepare("%h[fgRed]blue")
	exp, res := "\x1b[31mblue", p.Get(true)
	if res != exp {
		t.Errorf("Expected %q but result was %q", exp, res)
	}
	exp, res = "blue", p.Get(false)
	if res != exp {
		t.Errorf("Expected %q but result was %q", exp, res)
	}
}
