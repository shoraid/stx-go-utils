package asyncutil

import (
	"fmt"
	"runtime/debug"
)

// Result represents the result of an asynchronous computation.
type Result[T any] struct {
	Value T
	Err   error
}

// OnPanic is a global handler that will be called whenever a panic is recovered.
// You can assign this to send error to Sentry, log, metrics, etc.
var OnPanic func(err error)

// SafeGo runs a function asynchronously and recovers from panics.
// It returns a channel that yields the result (value and error).
func SafeGo[T any](fn func() (T, error)) <-chan Result[T] {
	ch := make(chan Result[T], 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				var zero T
				err := fmt.Errorf("panic recovered: %v\n%s", r, debug.Stack())

				if OnPanic != nil {
					// protect OnPanic from panicking
					defer func() {
						if rec := recover(); rec != nil {
							fmt.Printf("panic in OnPanic: %v\n", rec)
						}
					}()
					OnPanic(err)
				}

				ch <- Result[T]{Value: zero, Err: err}
			}
		}()

		val, err := fn()
		ch <- Result[T]{Value: val, Err: err}
	}()

	return ch
}
