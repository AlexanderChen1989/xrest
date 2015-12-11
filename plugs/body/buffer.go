package body

import "sync"

// Use simple []byte instead of bytes.Buffer to avoid large dependency.
type buffer []byte

func (b *buffer) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *buffer) free() {
	// Don't hold on to large buffers.
	if len(*b) <= 1024 {
		bufferFree.Put((*b)[:0])
	}
}

var bufferFree = &sync.Pool{
	New: func() interface{} {
		return buffer{}
	},
}

func getBuffer() buffer {
	return bufferFree.Get().(buffer)
}
