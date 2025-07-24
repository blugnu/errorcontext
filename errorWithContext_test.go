package errorcontext_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/blugnu/errorcontext"
)

func TestErrorWithContextError(t *testing.T) {
	// ARRANGE
	testcases := []struct {
		sut    error
		result string
	}{
		{sut: errorcontext.ErrorWithContext{}, result: "unspecified error with context"},
		{sut: errorcontext.Wrap(context.TODO(), errors.New("wrapped")), result: "wrapped"},
		{sut: errorcontext.Wrap(context.TODO(), fmt.Errorf("%w: %w", errors.New("wrapped"), errors.New("cause"))), result: "wrapped: cause"},
		{sut: errorcontext.Wrap(context.TODO(), errors.Join(errors.New("wrapped"), errors.New("cause"))), result: "wrapped\ncause"},
	}
	for _, tc := range testcases {
		t.Run(tc.result, func(t *testing.T) {
			// ACT
			got := tc.sut.Error()

			// ASSERT
			wanted := tc.result
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})
	}
}

func TestErrorsIs(t *testing.T) {
	// ARRANGE
	ea := errors.New("a")
	eb := errors.New("b")

	testcases := []struct {
		name   string
		sut    error
		target error
		result bool
	}{
		{name: "ErrorWithContext{}:ErrorWithContext{}", sut: errorcontext.ErrorWithContext{}, target: errorcontext.ErrorWithContext{}, result: true},
		{name: "ErrorWithContext{}:<some error>", sut: errorcontext.ErrorWithContext{}, target: errors.New(""), result: false},
		{name: "ErrorWithContext{error:a}:ErrorWithContext{error:nil}", sut: errorcontext.Wrap(context.TODO(), ea), target: errorcontext.ErrorWithContext{}, result: true},
		{name: "ErrorWithContext{error:a}:ErrorWithContext{error:a}", sut: errorcontext.Wrap(context.TODO(), ea), target: errorcontext.Wrap(context.TODO(), ea), result: true},
		{name: "ErrorWithContext{error:a}:ErrorWithContext{error:b}", sut: errorcontext.Wrap(context.TODO(), ea), target: errorcontext.Wrap(context.TODO(), eb), result: false},
		{name: "ErrorWithContext{error:a}:error(a)", sut: errorcontext.Wrap(context.TODO(), ea), target: ea, result: true},
		{name: "ErrorWithContext{error:a}:error(b)", sut: errorcontext.Wrap(context.TODO(), ea), target: eb, result: false},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ASSERT
			wanted := tc.result
			got := errors.Is(tc.sut, tc.target)
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})
	}
}

func TestErrorsAs(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	err := errors.New("error")

	testcases := []struct {
		name   string
		sut    error
		target interface{}
		result bool
	}{
		{name: "ErrorWithContext{}->ErrorWithContext{}", sut: errorcontext.Wrap(ctx, err), target: errorcontext.ErrorWithContext{}, result: true},
		{name: "ErrorWithContext{}->&ErrorWithContext{}", sut: errorcontext.Wrap(ctx, err), target: &errorcontext.ErrorWithContext{}, result: true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ASSERT
			wanted := tc.result
			got := errors.As(tc.sut, &tc.target)
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})
	}
}

func TestErrorWithContextContext(t *testing.T) {
	// ARRANGE
	type keytype int
	const key keytype = 1
	ctxa := context.Background()
	ctxb := context.WithValue(ctxa, key, "value")

	testcases := []struct {
		name   string
		sut    error
		result context.Context
	}{
		{name: "wraps error with context", sut: errorcontext.Wrap(ctxa, errorcontext.Wrap(ctxb, errors.New("error"))), result: ctxb},
		{name: "does not wrap error with context", sut: errorcontext.Wrap(ctxa, errors.New("error")), result: ctxa},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			got := tc.sut.(errorcontext.ErrorWithContext).Context()

			// ASSERT
			wanted := tc.result
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})
	}
}
