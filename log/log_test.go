package log

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/nhooyr/color"
)

func TestPrintf(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	l := New(&b, true)
	const s = "%h[fgBlue]bar:%r %s"
	f := color.Prepare(s)
	f2 := color.Prepare("%h[fgWhite]bar")
	exp := fmt.Sprintf(f.Get(true), f2.Get(true)) + "\n"
	l.Printf(s, f2)
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
	b.Reset()
	l.SetColor(false)
	exp = fmt.Sprintf(f.Get(false), f2.Get(false)) + "\n"
	l.Printf(s, f2)
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}

func TestPrintfp(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	l := New(&b, true)
	f := color.Prepare("%h[fgBlue]bar:%r %s\n")
	f2 := color.Prepare("%h[fgWhite]bar")
	exp := fmt.Sprintf(f.Get(true), f2.Get(true))
	l.Printfp(f, f2)
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
	b.Reset()
	l.SetColor(false)
	exp = fmt.Sprintf(f.Get(false), f2.Get(false))
	l.Printfp(f, f2)
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}

func TestPrint(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	l := New(&b, true)
	f := color.Prepare("%h[fgWhite]bar")
	exp := f.Get(true) + "foo\n"
	l.Print(f, "foo")
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
	b.Reset()
	exp = f.Get(false) + "foo\n"
	l.SetColor(false)
	l.Print(f, "foo")
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}

func TestPrintln(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	l := New(&b, true)
	f := color.Prepare("%h[fgWhite]bar")
	exp := f.Get(true) + " foo\n"
	l.Println(f, "foo")
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
	b.Reset()
	exp = f.Get(false) + " foo\n"
	l.SetColor(false)
	l.Println(f, "foo")
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}

func TestPanic(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	SetOutput(&b)
	exp := "foohi"
	defer func() {
		if r, ok := recover().(string); !ok || r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		} else if b.String() != exp+"\n" {
			t.Errorf("Expected %q but result was %q", exp+"\n", b.String())
		}
	}()
	Panic("foo", "hi")
	panic("Impossible")
}
