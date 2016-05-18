package log

import (
	"bytes"
	"testing"
)

func benchmarkPrintln(b *testing.B, msg string) {
	buf := bytes.NewBuffer(nil)
	l := New(buf, false)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		l.Println(msg)
	}
}

func BenchmarkPrintlnShort(b *testing.B) {
	benchmarkPrintln(b, "test")
}

func BenchmarkPrintlnLong(b *testing.B) {
	benchmarkPrintln(b, "this sentence is many bytes long very long much long wow")
}
