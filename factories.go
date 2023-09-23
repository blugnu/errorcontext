package errorcontext

import (
	"context"
	"errors"
	"fmt"
)

// New creates a new error with the supplied Context.
// The error wraps a new error created by passing the
// specified string to `errors.New()`.
func New(ctx context.Context, s string) ErrorWithContext {
	return ErrorWithContext{ctx, errors.New(s)}
}

// Errorf creates a new error with the supplied Context.
// The error wraps an error created by passing a supplied
// format string and args to `fmt.Errorf()`.
//
// Since `fmt.Errorf` is used, %w may be used in the format
// string and an existing error passed as appropriate.
func Errorf(ctx context.Context, format string, args ...interface{}) ErrorWithContext {
	return ErrorWithContext{ctx, fmt.Errorf(format, args...)}
}

// Wrap creates a new error wrapping the specified error
// with the provided context.
func Wrap(ctx context.Context, err error) ErrorWithContext {
	return ErrorWithContext{ctx, err}
}
