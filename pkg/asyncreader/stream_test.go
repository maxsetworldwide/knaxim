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
	"bytes"
	"crypto/rand"
	"io"
	"sync"
	"testing"
)

func TestStream(t *testing.T) {
	w, rs := NewWithMaxsize(5, 10)
	savedBuf := new(bytes.Buffer)
	Treader := io.TeeReader(rand.Reader, savedBuf)
	wg := &sync.WaitGroup{}
	wg.Add(6)
	go func() {
		defer wg.Done()
		if copied, err := io.CopyN(w, Treader, 32); err != nil {
			t.Errorf("unable to write data(%d): %s", copied, err)
		}
		w.Close()
	}()
	results := make([]*bytes.Buffer, 5)
	for i := 0; i < 5; i++ {
		results[i] = new(bytes.Buffer)
		go func(b *bytes.Buffer, r io.Reader, indx int) {
			defer wg.Done()
			if copied, err := io.Copy(b, r); err != nil {
				t.Errorf("failed to read %d's data(%d): %s", indx, copied, err)
			}
		}(results[i], rs[i], i)
	}
	wg.Wait()
	if t.Failed() {
		t.FailNow()
	}
	for i := 0; i < 5; i++ {
		if !bytes.Equal(savedBuf.Bytes(), results[i].Bytes()) {
			t.Errorf("incorrectly copied values: %d, expected: %v, resulted: %v", i, savedBuf.Bytes(), results[i].Bytes())
		}
	}
}
