package color

import (
	"testing"

	"github.com/nhooyr/terminfo/caps"
)

func TestPrepare(t *testing.T) {
	t.Parallel()
	f := Prepare("%h[fgBlue]foo")
	exp := ti.Color(caps.Blue, -1) + "foo"
	r := f.Get(true)
	if exp != r {
		t.Errorf("Expected %q but result was %q", exp, r)
	}
	exp = "foo"
	r = f.Get(false)
	if exp != r {
		t.Errorf("Expected %q but result was %q", exp, r)
	}
}

func TestEprintf(t *testing.T) {
	t.Parallel()
	f := Prepare("%h[fgRed]panic: %s: %s").Eprintfp("bar", Prepare("%h[fgGreen]rip"))
	exp := ti.Color(caps.Red, -1) + "panic: bar: " + ti.Color(caps.Green, -1) + "rip"
	r := f.Get(true)
	if exp != r {
		t.Errorf("Expected %q but result was %q", exp, r)
	}
	exp = "panic: bar: rip"
	r = f.Get(false)
	if exp != r {
		t.Errorf("Expected %q but result was %q", exp, r)
	}
}

func TestExpandFormats(t *testing.T) {
	t.Parallel()
	a := [3]interface{}{
		Prepare("%h[bgMagenta]foo"),
		"bar",
		3,
	}
	exp := a
	exp[0]= ti.Color(-1, caps.Magenta) + "foo"
	r := a
	ExpandFormats(true, r[:])
	if exp != r {
		t.Errorf("Expected %q but result was %q", exp, r)
	}
	exp[0] = "foo"
	r = a
	ExpandFormats(false, r[:])
	if exp != r {
		t.Errorf("Expected %q but result was %q", exp, r)
	}
}
