package errorcontext_test

import (
	"context"
	"errors"
	"testing"

	"github.com/blugnu/test"

	"github.com/blugnu/errorcontext"
)

func TestFromReturnsDefaultContextIfNoErrorWithContextIsFound(t *testing.T) {
	// ARRANGE
	defaultCtx := context.Background()
	err := errors.New("no context")

	// ACT

	result := errorcontext.From(defaultCtx, err)

	// ASSERT
	test.Value(t, result).Equals(defaultCtx)
}

func TestFromReturnsContextFromAnErrorContext(t *testing.T) {
	// ARRANGE
	type key string

	baseCtx := context.Background()
	valueCtx := context.WithValue(baseCtx, key("k"), "value")
	err := errorcontext.New(valueCtx, "")

	// ACT
	result := errorcontext.From(baseCtx, err)

	// ASSERT
	test.Value(t, result).Equals(valueCtx)
}

func TestFromReturnsTheDeepestContext(t *testing.T) {
	// ARRANGE
	type keytype int
	const key keytype = 1

	baseCtx := context.Background()
	valueCtx := context.WithValue(baseCtx, key, "value")

	err := errorcontext.Wrap(baseCtx, errorcontext.Wrap(valueCtx, errors.New("no context")))

	// ACT
	result := errorcontext.From(baseCtx, err)

	// ASSERT
	test.Value(t, result).Equals(valueCtx)
}
