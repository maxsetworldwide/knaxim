package asyncreader

import (
	"io"
	"sync"
)

// NewWithMaxsize builds buffered pipes allow multiple readers to read
// the full content written. intended for use in concurrent setting
// n is the number of readers to create, if n == 1 its recommended to use bytes.Buffer
// m is the maxsize of buffered space before blocking until all readers have read some space
// if m is zero, it is considered to have no maxsize
func NewWithMaxsize(n int, m int) (io.WriteCloser, []io.Reader) {
	w := new(buffer)
	w.maxsize = m
	w.lock = new(sync.RWMutex)
	rs := make([]io.Reader, n)
	w.readers = make([]*bufferReader, n)
	for i := 0; i < n; i++ {
		w.readers[i] = &bufferReader{
			buffer: w,
		}
		rs[i] = w.readers[i]
	}
	return w, rs
}

// New calls NewWithMaxsize with m as 0
func New(n int) (io.WriteCloser, []io.Reader) {
	return NewWithMaxsize(n, 0)
}
