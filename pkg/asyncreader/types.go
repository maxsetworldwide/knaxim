// Copyright August 2020 Maxset Worldwide Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	newData        *sync.Cond
	readers        []*bufferReader
	shifter        *sync.Once
	closed         bool
}

func (buf *buffer) Write(data []byte) (n int, err error) {
	if buf.closed {
		return 0, io.ErrClosedPipe
	}
	b := data
	buf.lock.Lock()
	defer buf.lock.Unlock()
	for buf.maxsize > 0 && buf.maxsize < len(buf.data)+len(b) {
		if buf.maxsize > len(buf.data) {
			shift := buf.maxsize - len(buf.data)
			buf.data = append(buf.data, b[0:shift]...)
			b = b[shift:]
			buf.newData.Broadcast()
		}
		if buf.availableSpace == nil {
			buf.availableSpace = sync.NewCond(buf.lock)
		}
		buf.availableSpace.Wait()
	}
	buf.data = append(buf.data, b...)
	buf.newData.Broadcast()
	return len(data), nil
}

func (buf *buffer) Close() error {
	if !buf.closed {
		buf.lock.Lock()
		defer buf.lock.Unlock()
		buf.closed = true
		buf.newData.Broadcast()
		buf.newData = nil
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
	for len(br.data)-br.head <= 0 && br.newData != nil {
		br.newData.Wait()
	}
	n = copy(b, br.data[br.head:])
	br.head += n
	if br.closed && br.head >= len(br.data) {
		err = io.EOF
	}
	br.lock.RUnlock()
	go br.shift()
	return
}
