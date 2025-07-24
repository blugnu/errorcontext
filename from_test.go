package errorcontext_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/errorcontext"
)

func TestFromReturnsDefaultContextIfNoErrorWithContextIsFound(t *testing.T) {
	With(t)

	// arrange
	defaultCtx := context.Background()
	err := errors.New("no context")

	// act

	result := errorcontext.From(defaultCtx, err)

	// assert
	Expect(result).To(Equal(defaultCtx))
}

func TestFromReturnsContextFromAnErrorContext(t *testing.T) {
	With(t)

	// arrange
	type key string

	baseCtx := context.Background()
	valueCtx := context.WithValue(baseCtx, key("k"), "value")
	err := errorcontext.New(valueCtx, "")

	// act
	result := errorcontext.From(baseCtx, err)

	// assert
	Expect(result).To(Equal(valueCtx))
}

func TestFromReturnsTheDeepestContext(t *testing.T) {
	With(t)

	// arrange
	type keytype int
	const key keytype = 1

	baseCtx := context.Background()
	valueCtx := context.WithValue(baseCtx, key, "value")

	err := errorcontext.Wrap(baseCtx, errorcontext.Wrap(valueCtx, errors.New("no context")))

	// act
	result := errorcontext.From(baseCtx, err)

	// assert
	Expect(result).To(Equal(valueCtx))
}
