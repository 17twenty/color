package log

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/nhooyr/color"
	"github.com/uber-common/zap"
)

func TestPrintfAndPrintfp(t *testing.T) {
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
	l.Printfp(f, f2)
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
	b.Reset()
	l.Printfp(f, f2)
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}

func TestPrintAndPrintln(t *testing.T) {
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
	exp = f.Get(true) + " foo\n"
	l.Println(f, "foo")
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
	b.Reset()
	exp = f.Get(false) + " foo\n"
	l.Println(f, "foo")
	if b.String() != exp {
		t.Errorf("Expected %q but result was %q", exp, b.String())
	}
}

func TestPanic(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	l := New(&b, false)
	exp := "foohi"
	defer func() {
		if r, ok := recover().(string); !ok || r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		} else if b.String() != exp+"\n" {
			t.Errorf("Expected %q but result was %q", exp+"\n", b.String())
		}
	}()
	l.Panic("foo", "hi")
	panic("Impossible")
}

func TestPanicln(t *testing.T) {
	t.Parallel()
	var b bytes.Buffer
	SetOutput(&b)
	exp := "foo hi\n"
	defer func() {
		if r, ok := recover().(string); !ok || r != exp {
			t.Errorf("Expected %q but result was %q", exp, r)
		} else if b.String() != exp {
			t.Errorf("Expected %q but result was %q", exp, b.String())
		}
	}()
	Panicln("foo", "hi")
	panic("Impossible")
}

func BenchmarkPrintln(b *testing.B) {
	l := New(ioutil.Discard, true)
	for i := 0; i < b.N; i++ {
		l.Println("LOL")
	}
}

func BenchmarkPrintlnParallel(b *testing.B) {
	l := New(ioutil.Discard, true)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Println("foo")
		}
	})
}

func BenchmarkZap(b *testing.B) {
	l := zap.NewJSON(zap.Output(zap.AddSync(ioutil.Discard)))
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("foo")
		}
	})
}
