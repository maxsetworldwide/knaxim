/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

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
	w.newData = sync.NewCond(w.lock.RLocker())
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
