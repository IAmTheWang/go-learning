package errorsinterfaces

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNotFoundErrorError(t *testing.T) {
	err := &NotFoundError{ID: 42}
	want := "item 42 not found"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestFindItem(t *testing.T) {
	items := map[int]string{1: "apple", 2: "banana"}

	t.Run("found", func(t *testing.T) {
		got, err := FindItem(items, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "apple" {
			t.Errorf("FindItem(1) = %q, want %q", got, "apple")
		}
	})

	t.Run("missing", func(t *testing.T) {
		_, err := FindItem(items, 99)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		var nf *NotFoundError
		if !errors.As(err, &nf) {
			t.Fatalf("error is not a *NotFoundError: %v", err)
		}
		if nf.ID != 99 {
			t.Errorf("NotFoundError.ID = %d, want %d", nf.ID, 99)
		}
	})
}

func TestLoadItemConfig(t *testing.T) {
	items := map[int]string{1: "apple"}

	if err := LoadItemConfig(items, 1); err != nil {
		t.Errorf("unexpected error for existing item: %v", err)
	}

	err := LoadItemConfig(items, 99)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	const wantSubstr = "load config"
	if got := err.Error(); !strings.Contains(got, wantSubstr) {
		t.Errorf("error message %q does not mention %q", got, wantSubstr)
	}

	// The original *NotFoundError must still be reachable through the
	// wrapped error chain — this is the entire point of using %w.
	var nf *NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("wrapped error does not unwrap to *NotFoundError: %v", err)
	}
	if nf.ID != 99 {
		t.Errorf("unwrapped NotFoundError.ID = %d, want %d", nf.ID, 99)
	}
}

func TestIsNotFound(t *testing.T) {
	items := map[int]string{1: "apple"}

	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"plain error", errors.New("boom"), false},
		{"direct NotFoundError", &NotFoundError{ID: 5}, true},
		{"wrapped via LoadItemConfig", LoadItemConfig(items, 99), true},
	}

	for _, c := range cases {
		if got := IsNotFound(c.err); got != c.want {
			t.Errorf("%s: IsNotFound(%v) = %v, want %v", c.name, c.err, got, c.want)
		}
	}
}

func TestTemperatureString(t *testing.T) {
	cases := []struct {
		temp Temperature
		want string
	}{
		{23.5, "23.5°C"},
		{0, "0.0°C"},
		{-10, "-10.0°C"},
	}

	for _, c := range cases {
		if got := c.temp.String(); got != c.want {
			t.Errorf("Temperature(%v).String() = %q, want %q", float64(c.temp), got, c.want)
		}
		// fmt should call String() automatically for %v — this is what
		// makes the Stringer interface useful, not just callable-by-hand.
		if got := fmt.Sprintf("%v", c.temp); got != c.want {
			t.Errorf("fmt.Sprintf(%%v, %v) = %q, want %q (Stringer not picked up by fmt)", float64(c.temp), got, c.want)
		}
	}
}
