package asyncreader

import (
	"io"
	"sync"
)

type buffer struct {
	data           []byte
	maxsize        int
	lock           *sync.RWMutex
	availableSpace *sync.Cond
	readers        []*bufferReader
	shifter        *sync.Once
	closed         bool
}

func (buf *buffer) Write(data []byte) (n int, err error) {
	b := data
	buf.lock.Lock()
	defer buf.lock.Unlock()
	for buf.maxsize > 0 && buf.maxsize < len(buf.data)+len(b) {
		if buf.maxsize > len(buf.data) {
			buf.data = append(buf.data, b[0:buf.maxsize-len(buf.data)]...)
			b = b[buf.maxsize-len(buf.data):]
		}
		if buf.availableSpace == nil {
			buf.availableSpace = sync.NewCond(buf.lock)
		}
		buf.availableSpace.Wait()
	}
	buf.data = append(buf.data, b...)
	return len(data), nil
}

func (buf *buffer) Close() error {
	if !buf.closed {
		buf.lock.Lock()
		defer buf.lock.Unlock()
		buf.closed = true
	}
	return nil
}

func (buf *buffer) shift() {
	if buf.shifter == nil {
		buf.lock.Lock()
		if buf.shifter == nil {
			buf.shifter = new(sync.Once)
		}
		buf.lock.Unlock()
	}
	buf.shifter.Do(func() {
		buf.lock.Lock()
		defer buf.lock.Unlock()
		minpos := -1
		for _, br := range buf.readers {
			if minpos == -1 || minpos > br.head {
				minpos = br.head
			}
		}
		if minpos > 0 {
			buf.data = buf.data[minpos:]
			for _, br := range buf.readers {
				br.head = br.head - minpos
			}
			if buf.availableSpace != nil {
				buf.availableSpace.Broadcast()
			}
		}
		buf.shifter = new(sync.Once)
	})
}

type bufferReader struct {
	*buffer
	head int
}

func (br *bufferReader) Read(b []byte) (n int, err error) {
	br.lock.RLock()
	n = copy(b, br.data[br.head:])
	br.head += n
	if br.closed && br.head >= len(br.data) {
		err = io.EOF
	}
	br.lock.RUnlock()
	go br.shift()
	return
}
