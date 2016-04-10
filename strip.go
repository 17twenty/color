package color

import (
	"strings"
	"sync"
)

// bufferPool allows the reuse of buffers to avoid allocations.
var bufferPool = sync.Pool{
	New: func() interface{} {
		// The initial capacity avoids constant reallocation during growth.
		return buffer(make([]byte, 0, 30))
	},
}

// Sstripf removes all highlight verbs in s and then returns the resulting string.
// This is a low level function, you shouldn't need to use this most of the time.
func Sstripf(s string) string {
	buf := bufferPool.Get().(buffer)
	// pi is the index after the last verb.
	var pi, i int
LOOP:
	for ; ; i++ {
		if i >= len(s) {
			if i > pi {
				buf.writeString(s[pi:i])
			}
			break
		} else if s[i] != '%' {
			continue
		}
		if i > pi {
			buf.writeString(s[pi:i])
		}
		i++
		if i >= len(s) {
			// Let fmt handle "%!h(NOVERB)".
			buf.writeByte('%')
			break
		}
		switch s[i] {
		case 'r':
			// Strip the reset verb.
			pi = i + 1
		case 'h':
			// Strip inside the highlight verb.
			j := strings.IndexByte(s[i+1:], ']')
			if j == -1 {
				buf.writeString(errInvalid)
				break LOOP
			}
			i += j + 1
			pi = i + 1
		default:
			// Include the verb.
			pi = i - 1
		}
	}
	s = string(buf)
	buf.reset()
	bufferPool.Put(buf)
	return s
}
