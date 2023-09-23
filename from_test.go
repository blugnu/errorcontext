package errorcontext

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestFromErrorReturnsDefaultContextIfNoErrorWithContextIsFound(t *testing.T) {

	// ARRANGE
	ctx := context.Background()
	err := errors.New("no context")

	// ACT

	got := From(ctx, err)

	// ASSERT
	wanted := ctx
	if wanted != got {
		t.Errorf("wanted %v, got %v", wanted, got)
	}
}

func TestFromErrorReturnsTheDeepestContext(t *testing.T) {

	// ARRANGE
	type keytype int
	const key keytype = 1

	initial := context.Background()
	ctx := context.WithValue(initial, key, "value")

	err := &ErrorWithContext{ctx: ctx, error: errors.New("no context")}
	buried := fmt.Errorf("buried: %w", err)
	inner := fmt.Errorf("inner: %w", buried)
	outer := fmt.Errorf("outer: %w", inner)

	// ACT
	got := From(initial, outer)

	// ASSERT
	wanted := ctx
	if wanted != got {
		t.Errorf("wanted %v, got %v", ctx, got)
	}
}
