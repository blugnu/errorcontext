package errorcontext

import (
	"context"
	"errors"
)

// ErrorWithContext wraps an error with a context.
type ErrorWithContext struct {
	ctx context.Context
	error
}

// Error implements the error interface.
func (err ErrorWithContext) Error() string {
	if err.error != nil {
		return err.error.Error()
	}
	return "unknown error"
}

// Context returns the inner-most context accessible from
// this error or any wrapped ErrorWithContext.
//
// That is, the wrapped error is tested to determine if it
// is (or wraps) an ErrorWithContext and if it is the Context()
// function on that wrapped error is called.  This continues
// recursively until there are no more ErrorWithContext errors
// to be unwrapped.
func (err ErrorWithContext) Context() context.Context {
	wrapped := ErrorWithContext{}
	if errors.As(err.error, &wrapped) {
		return wrapped.Context()
	}
	return err.ctx
}

// Is compares an ErrorWithContext with some target error to
// determine whether they are considered equal.
//
// A receiver will match with a target that:
// - is an ErrorWithContext or *ErrorWithContext, and
// - has an equal or nil context, and
// - has a nil error or
// - an error which satisfies errors.Is(receiver.error, target.error)
func (err ErrorWithContext) Is(target error) bool {
	match := func(target *ErrorWithContext) bool {
		return (target.ctx == nil || err.ctx == target.ctx) &&
			(target.error == nil || errors.Is(err.error, target.error))
	}

	switch target := target.(type) {
	case ErrorWithContext:
		return match(&target)
	case *ErrorWithContext:
		return match(target)
	default:
		return false
	}
}

// Unwrap returns the error wrapped by the ErrorWithContext.
func (err ErrorWithContext) Unwrap() error {
	return err.error
}
