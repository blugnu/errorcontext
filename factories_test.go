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
	sut := New(ctx, "some new error").(ErrorWithContext)

	// ASSERT
	assert(t, sut, ctx, func() bool { return sut.error != nil }, "some new error")
}

func Test_Errorf(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	err := errors.New("some error")

	// ACT
	sut := Errorf(ctx, "narrative: %w", err).(ErrorWithContext)

	// ASSERT
	assert(t, sut, ctx, func() bool { return errors.Is(sut, err) }, "narrative: some error")
}

func Test_Wrap(t *testing.T) {
	// ARRANGE
	ctx := context.Background()
	w := errors.New("wrapped")
	c := errors.New("cause")
	x := errors.New("extra")

	type result struct {
		string
	}
	testcases := []struct {
		name string
		args []error
		*result
	}{
		{name: "all nil (1)", args: []error{nil}, result: nil},
		{name: "all nil (2)", args: []error{nil, nil}, result: nil},
		{name: "all nil (3)", args: []error{nil, nil, nil}, result: nil},
		{name: "one non-nil (1)", args: []error{w}, result: &result{string: "wrapped"}},
		{name: "1st of 2 non-nil", args: []error{w, nil}, result: &result{string: "wrapped"}},
		{name: "2nd of 2 non-nil", args: []error{nil, w}, result: &result{string: "wrapped"}},
		{name: "both of 2 non-nil", args: []error{w, c}, result: &result{string: "wrapped: cause"}},
		{name: "1st of 3 non-nil", args: []error{w, nil, nil}, result: &result{string: "wrapped"}},
		{name: "2nd of 3 non-nil", args: []error{nil, w, nil}, result: &result{string: "wrapped"}},
		{name: "3rd of 3 non-nil", args: []error{nil, nil, w}, result: &result{string: "wrapped"}},
		{name: "1st and 2nd of 3 non-nil", args: []error{w, c, nil}, result: &result{string: "wrapped\ncause"}},
		{name: "1st and 3rd of 3 non-nil", args: []error{w, nil, c}, result: &result{string: "wrapped\ncause"}},
		{name: "2nd and 3rd of 3 non-nil", args: []error{nil, w, c}, result: &result{string: "wrapped\ncause"}},
		{name: "all of 3 non-nil", args: []error{w, c, x}, result: &result{string: "wrapped\ncause\nextra"}},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// ACT
			sut := Wrap(ctx, tc.args[0], tc.args[1:]...)

			// ASSERT
			switch {
			case tc.result == nil && sut == nil:
				return
			case tc.result == nil && sut != nil:
				t.Errorf("expected nil")
			case tc.result != nil && sut == nil:
				t.Errorf("expected error")
			default:
				for _, err := range tc.args {
					if err != nil && !errors.Is(sut, err) {
						t.Errorf("expected error to wrap: %v", err)
					}
				}
				wanted := tc.result.string
				got := sut.Error()
				if wanted != got {
					t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
				}
			}
		})
	}
}
