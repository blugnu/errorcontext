package errorcontext_test

import (
	"context"
	"errors"
	"testing"

	"github.com/blugnu/test"

	"github.com/blugnu/errorcontext"
)

func TestFactories(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	a := errors.New("a")
	b := errors.New("b")

	testcases := []struct {
		scenario string
		exec     func(t *testing.T)
	}{
		// Errorf tests
		{scenario: "Errorf(\"narrative: %w\",a)",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.Errorf(ctx, "narrative: %w", a)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("narrative: a")
				if result, ok := test.IsType[errorcontext.ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
				}
			},
		},

		// errorcontext.Join tests
		{scenario: "Join(nil)",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.Join(ctx, nil)

				// ASSERT
				test.IsNil(t, result)
			},
		},
		{scenario: "Join(a)",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.Join(ctx, a)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a")
				if result, ok := test.IsType[errorcontext.ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
				}
			},
		},
		{scenario: "Join(a,b)",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.Join(ctx, a, b)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a\nb")
				if result, ok := test.IsType[errorcontext.ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
					test.Error(t, result).Is(b)
				}
			},
		},

		// New tests
		{scenario: "New(ctx, \"a\")",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.New(ctx, "a")

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a")
				_, _ = test.IsType[errorcontext.ErrorWithContext](t, result)
			},
		},

		// errorcontext.Wrap tests
		{scenario: "Wrap(nil)",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.Wrap(ctx, nil)

				// ASSERT
				test.IsNil(t, result)
			},
		},
		{scenario: "Wrap(a)",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.Wrap(ctx, a)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a")
				if result, ok := test.IsType[errorcontext.ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
				}
			},
		},
		{scenario: "Wrap(a, nil)",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.Wrap(ctx, a, nil)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a")
				if result, ok := test.IsType[errorcontext.ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
				}
			},
		},
		{scenario: "Wrap(nil,b)",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.Wrap(ctx, nil, b)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("b")
				if result, ok := test.IsType[errorcontext.ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(b)
				}
			},
		},
		{scenario: "Wrap(a,b)",
			exec: func(t *testing.T) {
				// ACT
				result := errorcontext.Wrap(ctx, a, b)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a: b")
				if result, ok := test.IsType[errorcontext.ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
					test.Error(t, result).Is(b)
				}
			},
		},
		{scenario: "Wrap(nil,nil,nil)",
			exec: func(t *testing.T) {
				// ARRANGE ASSERT
				defer test.ExpectPanic(errorcontext.ErrIllegalOperation).Assert(t)

				// ACT
				_ = errorcontext.Wrap(ctx, nil, nil, nil)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
