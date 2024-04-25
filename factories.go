package errorcontext

import (
	"context"
	"errors"
	"fmt"
)

var ErrIllegalOperation = errors.New("illegal operation")

// Errorf creates a new error with the supplied Context.
// The error wraps an error created by passing a supplied
// format string and args to `fmt.Errorf()`.
//
// Since `fmt.Errorf` is used, %w may be used in the format
// string and an existing error passed as appropriate.
func Errorf(ctx context.Context, format string, args ...interface{}) error {
	return ErrorWithContext{ctx, fmt.Errorf(format, args...)}
}

// Join creates a new error wrapping any number of supplied errors
// with a specified context.
//
// If all errors are nil then the function returns nil.
func Join(ctx context.Context, errs ...error) error {
	if err := errors.Join(errs...); err != nil {
		return ErrorWithContext{ctx, err}
	}
	return nil
}

// New creates a new error with the supplied Context.
// The error wraps a new error created by passing the
// specified string to `errors.New()`.
func New(ctx context.Context, s string) error {
	return ErrorWithContext{ctx, errors.New(s)}
}

// Wrap creates a new error wrapping one or two errors with a
// specified context.  If both errors are nil then the function
// returns nil.
//
// The error returned by the function wraps the specified error(s)
// differently, depending on the number of non-nil errors specified.
//
// # Called with 1 Error
//
// If only one non-nil error is specified it is wrapped directly.
//
// # Called with 2 Errors
//
// If two errors (a, b) are specified and both are not nil, then both
// errors are wrapped within a new error of the form "a: b".
//
// i.e. the following are equivalent statements when a and b are non-nil:
//
//	errorcontext.Wrap(ctx, a, b)
//	errorcontext.Wrap(ctx, fmt.Errorf("%w: %w", a, b))
//
// If either a or b are nil (but not both), the non-nil error is
// wrapped directly.
//
// # Called with > 2 errors
//
// The function will panic with ErrIllegalOperation if called
// with 3 or more error arguments, whether or not any or all of
// those arguments are nil.  To wrap three or more errors use
// the errorcontext.Join function.
func Wrap(ctx context.Context, err error, errs ...error) error {
	switch 1 + len(errs) {
	case 1:
		if err != nil {
			return ErrorWithContext{ctx, err}
		}

	case 2:
		switch {
		case err == nil && errs[0] != nil:
			return ErrorWithContext{ctx, errs[0]}

		case err != nil && errs[0] == nil:
			return ErrorWithContext{ctx, err}

		case err != nil && errs[0] != nil:
			return ErrorWithContext{ctx, fmt.Errorf("%w: %w", err, errs[0])}
		}

	default:
		panic(fmt.Errorf("errorcontext.Wrap: %w: may only be called with 1 or 2 errors; "+
			"use errorcontext.Join for > 2 error arguments", ErrIllegalOperation))
	}
	return nil
}
