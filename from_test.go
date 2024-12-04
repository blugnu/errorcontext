package errorcontext

import (
	"context"
	"errors"
	"testing"

	"github.com/blugnu/test"
)

func TestFromReturnsDefaultContextIfNoErrorWithContextIsFound(t *testing.T) {
	// ARRANGE
	defaultCtx := context.Background()
	err := errors.New("no context")

	// ACT

	result := From(defaultCtx, err)

	// ASSERT
	test.Value(t, result).Equals(defaultCtx)
}

func TestFromReturnsContextFromAnErrorContext(t *testing.T) {
	// ARRANGE
	type key string

	baseCtx := context.Background()
	valueCtx := context.WithValue(baseCtx, key("k"), "value")
	err := ErrorWithContext{ctx: valueCtx}

	// ACT
	result := From(baseCtx, err)

	// ASSERT
	test.Value(t, result).Equals(valueCtx)
}

func TestFromReturnsTheDeepestContext(t *testing.T) {
	// ARRANGE
	type keytype int
	const key keytype = 1

	baseCtx := context.Background()
	valueCtx := context.WithValue(baseCtx, key, "value")

	err := ErrorWithContext{
		ctx: baseCtx,
		error: ErrorWithContext{
			ctx:   valueCtx,
			error: errors.New("no context"),
		},
	}

	// ACT
	result := From(baseCtx, err)

	// ASSERT
	test.Value(t, result).Equals(valueCtx)
}
