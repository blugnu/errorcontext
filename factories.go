package errorcontext

import (
	"context"
	"errors"
	"fmt"
)

// New creates a new error with the supplied Context.
// The error wraps a new error created by passing the
// specified string to `errors.New()`.
func New(ctx context.Context, s string) error {
	return ErrorWithContext{ctx, errors.New(s)}
}

// Errorf creates a new error with the supplied Context.
// The error wraps an error created by passing a supplied
// format string and args to `fmt.Errorf()`.
//
// Since `fmt.Errorf` is used, %w may be used in the format
// string and an existing error passed as appropriate.
func Errorf(ctx context.Context, format string, args ...interface{}) error {
	return ErrorWithContext{ctx, fmt.Errorf(format, args...)}
}

// Wrap creates a new error wrapping at least one error with a
// specified context.  However many arguments are specified,
// if all errors are nil then the function returns nil.
//
// The error returned by the function wraps the specified error(s)
// differently, depending on the number of non-nil errors specified.
//
// * 1 Error
//
// If only one non-nil error is specified, it is wrapped directly.
//
// * 2 Errors
//
// If two errors are specified and both are not nil, then the first
// is wrapped formatted to indicate the second as the cause.  i.e. both
// of the following are equivalent statements when a and b are non-nil:
//
//	errorcontext.Wrap(ctx, a, b)
//	errorcontext.Wrap(ctx, fmt.Errorf("%w: %w", a, b))
//
// If either a or b are nil (but not both), the non-nil error is
// wrapped directly.
//
// * 3 or More
//
// If three or more errors are specified, the function is
// equivalent to:
//
//	errorcontext.Wrap(ctx, errors.Join(errs...))
//
// Unless all errors are nil, in which case nil is returned.
func Wrap(ctx context.Context, err error, errs ...error) error {
	switch len(errs) {
	case 0:
		if err == nil {
			return nil
		}
		return ErrorWithContext{ctx, err}

	case 1:
		if err == nil && errs[0] == nil {
			return nil
		}
		if err != nil && errs[0] != nil {
			return ErrorWithContext{ctx, fmt.Errorf("%w: %w", err, errs[0])}
		}
		if err == nil {
			err = errs[0]
		}
		return ErrorWithContext{ctx, err}

	default:
		if err = errors.Join(append([]error{err}, errs...)...); err != nil {
			return ErrorWithContext{ctx, err}
		}
		return nil
	}
}
