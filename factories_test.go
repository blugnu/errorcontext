package errorcontext

import (
	"context"
	"errors"
	"testing"
)

func assert(t *testing.T, sut ErrorWithContext, wantCtx context.Context, errorIsWanted func() bool, wantString string) {
	t.Run("wraps context", func(t *testing.T) {
		wanted := wantCtx
		got := sut.ctx
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})

	t.Run("wraps error", func(t *testing.T) {
		wanted := true
		got := errorIsWanted()
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})

	t.Run("Error()", func(t *testing.T) {
		wanted := wantString
		got := sut.Error()
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})
}

func Test_New(t *testing.T) {
	// ARRANGE
	ctx := context.Background()

	// ACT
	sut := New(ctx, "some new error")

	// ASSERT
	assert(t, sut, ctx, func() bool { return sut.error != nil }, "some new error")
}

func Test_Errorf(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	err := errors.New("some error")

	// ACT
	sut := Errorf(ctx, "narrative: %w", err)

	// ASSERT
	assert(t, sut, ctx, func() bool { return errors.Is(sut, err) }, "narrative: some error")
}

func Test_Wrap(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	err := errors.New("some error")

	// ACT
	sut := Wrap(ctx, err)

	// ASSERT
	assert(t, sut, ctx, func() bool { return sut.error == err }, "some error")
}
