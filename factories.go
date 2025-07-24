package errorcontext

import (
	"context"
	"errors"
	"fmt"
)

// Errorf creates a new error with the supplied Context. The function
// is an analog for [fmt.Errorf].
func Errorf(ctx context.Context, format string, args ...any) error {
	return ErrorWithContext{ctx, fmt.Errorf(format, args...)}
}

// Join creates a new error wrapping any number of supplied errors
// with a specified context. The function is an analog for [errors.Join].
//
// If all errors are nil then the function returns nil.
func Join(ctx context.Context, errs ...error) error {
	if err := errors.Join(errs...); err != nil {
		return ErrorWithContext{ctx, err}
	}
	return nil
}

// New creates a new error with the supplied Context. The function is
// an analog for [errors.New].
func New(ctx context.Context, s string) error {
	return ErrorWithContext{ctx, errors.New(s)}
}

// Wrap creates a new error wrapping one or more errors with a
// specified context.
//
// The error returned by the function wraps the specified error(s)
// differently, depending on the number of non-nil errors specified.
//
// # Called with 1 Error
//
// If only one non-nil error is specified it is wrapped.
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
// wrapped.
//
// If both a and b are nil, the function returns nil.
//
// # Called with > 2 errors
//
// If more than two errors are specified, the function behaves
// similarly to the case with two errors, with all errors other than
// the first joined.
//
// i.e. the following are equivalent statements when a, b, and c are non-nil:
//
//	errorcontext.Wrap(ctx, a, b, c)
//	errorcontext.Wrap(ctx, fmt.Errorf("%w: %w", a, errors.Join(b, c)))
//
// If a is nil, the function returns an ErrorWithContext that wraps
// the joined errors b and c.
func Wrap(ctx context.Context, err error, errs ...error) error {
	switch len(errs) {
	case 0:
		// just err

	case 1:
		if err == nil {
			err = errs[0]
		} else if errs[0] != nil {
			err = fmt.Errorf("%w: %w", err, errs[0])
		}

	default: // err + 2 or more errors
		joined := errors.Join(errs...)
		if err == nil {
			err = joined
		} else {
			err = fmt.Errorf("%w: %w", err, joined)
		}
	}

	if err == nil {
		return nil
	}

	return ErrorWithContext{ctx, err}
}
