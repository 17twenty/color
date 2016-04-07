package color

type buffer []byte

func (b *buffer) write(p []byte) {
	*b = append(*b, p...)
}

func (b *buffer) writeString(s string) {
	*b = append(*b, s...)
}

func (b *buffer) reset() {
	*b = (*b)[:0]
}

func (b *buffer) writeByte(c byte) {
	*b = append(*b, c)
}
