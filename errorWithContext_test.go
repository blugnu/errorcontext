package errorcontext

import (
	"context"
	"errors"
	"testing"
)

func TestErrorWithContextError(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	err := errors.New("error")
	sut := &ErrorWithContext{ctx, err}

	// ACT
	s := sut.Error()

	// ASSERT
	wanted := "error"
	got := s
	if wanted != got {
		t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
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
		sut    *ErrorWithContext
		result context.Context
	}{
		{name: "wraps error with context", sut: &ErrorWithContext{ctxa, &ErrorWithContext{ctxb, errors.New("error")}}, result: ctxb},
		{name: "does not wrap error with context", sut: &ErrorWithContext{ctxa, errors.New("error")}, result: ctxa},
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
