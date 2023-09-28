package errorcontext

import (
	"context"
	"errors"
	"testing"
)

func TestErrorWithContextError(t *testing.T) {
	// ARRANGE
	// ARRANGE
	testcases := []struct {
		sut    ErrorWithContext
		result string
	}{
		{sut: ErrorWithContext{}, result: "unknown error"},
		{sut: ErrorWithContext{error: errors.New("wrapped error")}, result: "wrapped error"},
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
		{name: "ErrorWithContext{}:ErrorWithContext{}", sut: ErrorWithContext{}, target: ErrorWithContext{}, result: true},
		{name: "ErrorWithContext{}:&ErrorWithContext{}", sut: ErrorWithContext{}, target: &ErrorWithContext{}, result: true},
		{name: "&ErrorWithContext{}:ErrorWithContext{}", sut: &ErrorWithContext{}, target: ErrorWithContext{}, result: true},
		{name: "&ErrorWithContext{}:&ErrorWithContext{}", sut: &ErrorWithContext{}, target: &ErrorWithContext{}, result: true},
		{name: "ErrorWithContext{}:<some error>", sut: ErrorWithContext{}, target: errors.New(""), result: false},
		{name: "&ErrorWithContext{}:<some error>", sut: &ErrorWithContext{}, target: errors.New(""), result: false},
		{name: "ErrorWithContext{error:a}:ErrorWithContext{error:nil}", sut: ErrorWithContext{error: ea}, target: ErrorWithContext{}, result: true},
		{name: "ErrorWithContext{error:a}:ErrorWithContext{error:a}", sut: ErrorWithContext{error: ea}, target: ErrorWithContext{error: ea}, result: true},
		{name: "ErrorWithContext{error:a}:ErrorWithContext{error:b}", sut: ErrorWithContext{error: ea}, target: ErrorWithContext{error: eb}, result: false},
		{name: "ErrorWithContext{error:a}:error(a)", sut: ErrorWithContext{error: ea}, target: ea, result: true},
		{name: "ErrorWithContext{error:a}:error(b)", sut: ErrorWithContext{error: ea}, target: eb, result: false},
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
		{name: "ErrorWithContext{}->ErrorWithContext{}", sut: ErrorWithContext{ctx, err}, target: ErrorWithContext{}, result: true},
		{name: "ErrorWithContext{}->&ErrorWithContext{}", sut: ErrorWithContext{ctx, err}, target: &ErrorWithContext{}, result: true},
		{name: "&ErrorWithContext{}->ErrorWithContext{}", sut: &ErrorWithContext{ctx, err}, target: ErrorWithContext{}, result: true},
		{name: "&ErrorWithContext{}->&ErrorWithContext{}", sut: &ErrorWithContext{ctx, err}, target: &ErrorWithContext{}, result: true},
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
		sut    ErrorWithContext
		result context.Context
	}{
		{name: "wraps error with context", sut: ErrorWithContext{ctxa, ErrorWithContext{ctxb, errors.New("error")}}, result: ctxb},
		{name: "does not wrap error with context", sut: ErrorWithContext{ctxa, errors.New("error")}, result: ctxa},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			got := tc.sut.Context()

			// ASSERT
			wanted := tc.result
			if wanted != got {
				t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
			}
		})
	}
}
