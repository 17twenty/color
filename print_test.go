package color

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPrintf(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	const s = "%h[fgBlue]bar:%r %s\n"
	f := Prepare(s)
	f2 := Prepare("%h[fgWhite]bar")
	exp := fmt.Sprintf(f.Get(true), f2.Get(true))
	p := New(&b, true)
	p.Printf(s, f2)
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
	b.Reset()
	exp = fmt.Sprintf(f.Get(false), f2.Get(false))
	p = New(&b, false)
	p.Printf(s, f2)
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}

func TestPrintfp(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	f := Prepare("%h[fgBlue]bar:%r %s\n")
	f2 := Prepare("%h[fgWhite]bar")
	exp := fmt.Sprintf(f.Get(true), f2.Get(true))
	p := New(&b, true)
	p.Printfp(f, f2)
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
	b.Reset()
	exp = fmt.Sprintf(f.Get(false), f2.Get(false))
	p = New(&b, false)
	p.Printfp(f, f2)
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}

func TestPrint(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	f := Prepare("%h[fgWhite]bar")
	exp := f.Get(true) + "foo"
	p := New(&b, true)
	p.Print(f, "foo")
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
	b.Reset()
	exp = f.Get(false) + "foo"
	p = New(&b, false)
	p.Print(f, "foo")
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}

func TestPrintln(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	f := Prepare("%h[fgWhite]bar")
	exp := f.Get(true) + " foo\n"
	p := New(&b, true)
	p.Println(f, "foo")
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
	b.Reset()
	exp = f.Get(false) + " foo\n"
	p = New(&b, false)
	p.Println(f, "foo")
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}
