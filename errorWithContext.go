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
	return err.error.Error()
}

// Context returns the innermost context accessible from
// this error or any wrapped ErrorWithContext.
func (err ErrorWithContext) Context() context.Context {
	inner := &ErrorWithContext{}
	if errors.As(err.error, &inner) {
		return inner.Context()
	}
	return err.ctx
}

// Unwrap implements unwrapping to return the wrapped error.
func (err ErrorWithContext) Unwrap() error {
	return err.error
}
