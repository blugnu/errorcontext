package errorcontext_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/errorcontext"
)

func TestFactories(t *testing.T) {
	With(t)

	// ARRANGE
	ctx := context.Background()
	a := errors.New("a")
	b := errors.New("b")
	c := errors.New("c")

	Run(Test("Errorf(\"narrative: %w\",a)", func() {
		// act
		result := errorcontext.Errorf(ctx, "narrative: %w", a)

		// assert
		Expect(result).IsNotNil()
		Expect(result.Error()).To(Equal("narrative: a"))

		var err error = RequireType[errorcontext.ErrorWithContext](result)
		Expect(err).Is(a)
	}))

	// errorcontext.Join tests
	Run(Test("Join(nil)", func() {
		// ACT
		result := errorcontext.Join(ctx, nil)

		// ASSERT
		Expect(result).IsNil()
	}))

	Run(Test("Join(a)", func() {
		// ACT
		result := errorcontext.Join(ctx, a)

		// ASSERT
		Expect(result).IsNotNil()
		Expect(result.Error()).To(Equal("a"))
		var err error = RequireType[errorcontext.ErrorWithContext](result)
		Expect(err).Is(a)
	}))

	Run(Test("Join(a,b)", func() {
		// ACT
		result := errorcontext.Join(ctx, a, b)

		// ASSERT
		Expect(result).IsNotNil()
		Expect(result.Error()).To(Equal("a\nb"))

		var err error = RequireType[errorcontext.ErrorWithContext](result)
		Expect(err).Is(a)
	}))

	Run(Test("New(ctx, \"a\")", func() {
		// ACT
		result := errorcontext.New(ctx, "a")

		// ASSERT
		Expect(result).IsNotNil()
		Expect(result.Error()).To(Equal("a"))
		_ = RequireType[errorcontext.ErrorWithContext](result)
	}))

	Run(Test("Wrap(nil)", func() {
		// ACT
		result := errorcontext.Wrap(ctx, nil)

		// ASSERT
		Expect(result).IsNil()
	}))

	Run(Test("Wrap(a)", func() {
		// ACT
		result := errorcontext.Wrap(ctx, a)

		// ASSERT
		Expect(result).IsNotNil()
		Expect(result.Error()).To(Equal("a"))

		var err error = RequireType[errorcontext.ErrorWithContext](result)
		Expect(err).Is(a)
	}))

	Run(Test("Wrap(a, nil)", func() {
		// ACT
		result := errorcontext.Wrap(ctx, a, nil)

		// ASSERT
		Expect(result).IsNotNil()
		Expect(result.Error()).To(Equal("a"))

		var err error = RequireType[errorcontext.ErrorWithContext](result)
		Expect(err).Is(a)
	}))

	Run(Test("Wrap(nil,b)", func() {
		// ACT
		result := errorcontext.Wrap(ctx, nil, b)

		// ASSERT
		Expect(result).IsNotNil()
		Expect(result.Error()).To(Equal("b"))

		var err error = RequireType[errorcontext.ErrorWithContext](result)
		Expect(err).Is(b)
	}))

	Run(Test("Wrap(a,b)", func() {
		// ACT
		result := errorcontext.Wrap(ctx, a, b)

		// ASSERT
		Expect(result).IsNotNil()
		Expect(result.Error()).To(Equal("a: b"))

		var err error = RequireType[errorcontext.ErrorWithContext](result)
		Expect(err).Is(a)
		Expect(result).Is(b)
	}))

	Run(Test("Wrap(nil,nil,nil)", func() {
		err := errorcontext.Wrap(ctx, nil, nil, nil)
		Expect(err).IsNil()
	}))

	Run(Test("Wrap(a,b,c)", func() {
		err := errorcontext.Wrap(ctx, a, b, c)
		Expect(err).Is(a)
		Expect(err).Is(b)
		Expect(err).Is(c)
	}))
}
