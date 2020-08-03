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

package decode

import (
	"context"
	"sync"
	"time"
)

// ContextKey is a type to differentiate values from this package from other values when used as a key in a context key-value pair
type ContextKey byte

const (
	// PROCESSING is a key for a context value that should either be a buffered channel of struct{} or sync.Locker. These objects are used to limit the number of active processing threads. If unset processing run as soon as called.
	PROCESSING ContextKey = 'p'
	// TIMEOUT is a key for a context value that is expected to be a time.Duration. It is the time allotted to the processing of a file, if nil no time limit
	TIMEOUT       ContextKey = 't'
	timeoutCancel ContextKey = 'c'
)

func startProcessing(ctx context.Context) context.Context {
	if processinglock := ctx.Value(PROCESSING); processinglock != nil {
		switch pl := processinglock.(type) {
		case chan struct{}:
			<-pl
		case sync.Locker:
			pl.Lock()
		}
	}
	// add timeout to context after processing can begin
	if timeout := ctx.Value(TIMEOUT); timeout != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout.(time.Duration))
		ctx = context.WithValue(ctx, timeoutCancel, cancel)
	}
	return ctx
}

func stopProcessing(ctx context.Context) {
	if processinglock := ctx.Value(PROCESSING); processinglock != nil {
		switch pl := processinglock.(type) {
		case chan struct{}:
			pl <- struct{}{}
		case sync.Locker:
			pl.Unlock()
		}
	}
	if canceltimeout := ctx.Value(timeoutCancel); canceltimeout != nil {
		canceltimeout.(context.CancelFunc)()
	}
}
