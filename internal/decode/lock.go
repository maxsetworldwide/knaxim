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
