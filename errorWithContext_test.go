package errorcontext_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	. "github.com/blugnu/test"

	"github.com/blugnu/errorcontext"
)

func TestErrorWithContextError(t *testing.T) {
	With(t)

	type testcase struct {
		sut    error
		result string
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			got := tc.sut.Error()
			Expect(got).To(Equal(tc.result))
		}),
		Cases([]testcase{
			{sut: errorcontext.ErrorWithContext{}, result: "unspecified error with context"},
			{sut: errorcontext.Wrap(context.TODO(), errors.New("wrapped")), result: "wrapped"},
			{sut: errorcontext.Wrap(context.TODO(), fmt.Errorf("%w: %w", errors.New("wrapped"), errors.New("cause"))), result: "wrapped: cause"},
			{sut: errorcontext.Wrap(context.TODO(), errors.Join(errors.New("wrapped"), errors.New("cause"))), result: "wrapped\ncause"},
		}),
	))
}

func TestErrorsIs(t *testing.T) {
	With(t)

	// arrange
	ea := errors.New("a")
	eb := errors.New("b")

	type testcase struct {
		name   string
		sut    error
		target error
		result bool
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			got := errors.Is(tc.sut, tc.target)
			Expect(got).To(Equal(tc.result))
		}),
		Cases([]testcase{
			{name: "ErrorWithContext{}:ErrorWithContext{}", sut: errorcontext.ErrorWithContext{}, target: errorcontext.ErrorWithContext{}, result: true},
			{name: "ErrorWithContext{}:<some error>", sut: errorcontext.ErrorWithContext{}, target: errors.New(""), result: false},
			{name: "ErrorWithContext{error:a}:ErrorWithContext{error:nil}", sut: errorcontext.Wrap(context.TODO(), ea), target: errorcontext.ErrorWithContext{}, result: true},
			{name: "ErrorWithContext{error:a}:ErrorWithContext{error:a}", sut: errorcontext.Wrap(context.TODO(), ea), target: errorcontext.Wrap(context.TODO(), ea), result: true},
			{name: "ErrorWithContext{error:a}:ErrorWithContext{error:b}", sut: errorcontext.Wrap(context.TODO(), ea), target: errorcontext.Wrap(context.TODO(), eb), result: false},
			{name: "ErrorWithContext{error:a}:error(a)", sut: errorcontext.Wrap(context.TODO(), ea), target: ea, result: true},
			{name: "ErrorWithContext{error:a}:error(b)", sut: errorcontext.Wrap(context.TODO(), ea), target: eb, result: false},
		}),
	))
}

func TestErrorsAs(t *testing.T) {
	With(t)

	// arrange
	ctx := context.Background()
	err := errors.New("error")

	type testcase struct {
		name   string
		sut    error
		target any
		result bool
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			got := errors.As(tc.sut, &tc.target)
			Expect(got).To(Equal(tc.result))
		}),
		Cases([]testcase{
			{name: "ErrorWithContext{}->ErrorWithContext{}", sut: errorcontext.Wrap(ctx, err), target: errorcontext.ErrorWithContext{}, result: true},
			{name: "ErrorWithContext{}->&ErrorWithContext{}", sut: errorcontext.Wrap(ctx, err), target: &errorcontext.ErrorWithContext{}, result: true},
		}),
	))
}

func TestErrorWithContextContext(t *testing.T) {
	With(t)

	// arrange
	type keytype int
	const key keytype = 1
	ctxa := context.Background()
	ctxb := context.WithValue(ctxa, key, "value")

	type testcase struct {
		name   string
		sut    error
		result context.Context
	}
	Run(Testcases(
		ForEach(func(tc testcase) {
			got := tc.sut.(errorcontext.ErrorWithContext).Context()
			Expect(got).To(Equal(tc.result))
		}),
		Cases([]testcase{
			{name: "wraps error with context", sut: errorcontext.Wrap(ctxa, errorcontext.Wrap(ctxb, errors.New("error"))), result: ctxb},
			{name: "does not wrap error with context", sut: errorcontext.Wrap(ctxa, errors.New("error")), result: ctxa},
		}),
	))
}
