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

package brand

import (
	"math/rand"
	"sync"
	"time"
)

// ByteGen generates each possible byte in a random order
type ByteGen struct {
	r   *rand.Rand
	buf []byte
	i   int
	m   sync.Mutex
}

func (b *ByteGen) shuffle() {
	b.r.Shuffle(len(b.buf), b.getSwap())
	b.i = 0
}

func (b *ByteGen) getSwap() func(i, j int) {
	return func(i, j int) {
		b.buf[i], b.buf[j] = b.buf[j], b.buf[i]
	}
}

var bg = New(time.Now().Unix())

// New returns a new ByteGen using the given seed
func New(seed int64) *ByteGen {
	n := new(ByteGen)

	n.r = rand.New(rand.NewSource(seed))

	n.buf = make([]byte, 256)
	for i := range n.buf {
		n.buf[i] = byte(i)
	}

	n.shuffle()
	return n
}

// Next returns the next byte
func (b *ByteGen) Next() byte {
	b.m.Lock()
	defer b.m.Unlock()
	defer func() {
		b.i++
		if b.i >= len(bg.buf) {
			b.shuffle()
		}
	}()
	return b.buf[bg.i]
}

// Next returns a random byte without repeating until every possible byte has
// been returned
func Next() byte {
	return bg.Next()
}
