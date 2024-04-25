package errorcontext

import (
	"context"
	"errors"
	"testing"

	"github.com/blugnu/test"
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
				result := Errorf(ctx, "narrative: %w", a)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("narrative: a")
				if result, ok := test.IsType[ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
				}
			},
		},

		// Join tests
		{scenario: "Join(nil)",
			exec: func(t *testing.T) {
				// ACT
				result := Join(ctx, nil)

				// ASSERT
				test.IsNil(t, result)
			},
		},
		{scenario: "Join(a)",
			exec: func(t *testing.T) {
				// ACT
				result := Join(ctx, a)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a")
				if result, ok := test.IsType[ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
				}
			},
		},
		{scenario: "Join(a,b)",
			exec: func(t *testing.T) {
				// ACT
				result := Join(ctx, a, b)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a\nb")
				if result, ok := test.IsType[ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
					test.Error(t, result).Is(b)
				}
			},
		},

		// New tests
		{scenario: "New(ctx, \"a\")",
			exec: func(t *testing.T) {
				// ACT
				result := New(ctx, "a")

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a")
				_, _ = test.IsType[ErrorWithContext](t, result)
			},
		},

		// Wrap tests
		{scenario: "Wrap(nil)",
			exec: func(t *testing.T) {
				// ACT
				result := Wrap(ctx, nil)

				// ASSERT
				test.IsNil(t, result)
			},
		},
		{scenario: "Wrap(a)",
			exec: func(t *testing.T) {
				// ACT
				result := Wrap(ctx, a)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a")
				if result, ok := test.IsType[ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
				}
			},
		},
		{scenario: "Wrap(a, nil)",
			exec: func(t *testing.T) {
				// ACT
				result := Wrap(ctx, a, nil)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a")
				if result, ok := test.IsType[ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
				}
			},
		},
		{scenario: "Wrap(nil,b)",
			exec: func(t *testing.T) {
				// ACT
				result := Wrap(ctx, nil, b)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("b")
				if result, ok := test.IsType[ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(b)
				}
			},
		},
		{scenario: "Wrap(a,b)",
			exec: func(t *testing.T) {
				// ACT
				result := Wrap(ctx, a, b)

				// ASSERT
				test.IsTrue(t, result != nil)
				test.That(t, result.Error()).Equals("a: b")
				if result, ok := test.IsType[ErrorWithContext](t, result); ok {
					test.Error(t, result).Is(a)
					test.Error(t, result).Is(b)
				}
			},
		},
		{scenario: "Wrap(nil,nil,nil)",
			exec: func(t *testing.T) {
				// ARRANGE ASSERT
				defer test.ExpectPanic(ErrIllegalOperation).Assert(t)

				// ACT
				_ = Wrap(ctx, nil, nil, nil)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.scenario, func(t *testing.T) {
			tc.exec(t)
		})
	}
}
