package errorcontext

import (
	"context"
	"errors"
)

// From accepts a current context and an error and returns the context
// from the 'most wrapped' error with Context.  If the error is not an
// ErrorWithContext (and does not wrap one) the supplied `context` is
// returned.
func From(ctx context.Context, err error) context.Context {
	wrapped := ErrorWithContext{}
	if errors.As(err, &wrapped) {
		return wrapped.Context()
	}
	return ctx
}
