package functions

import "testing"

func TestDivide(t *testing.T) {
	got, err := Divide(10, 2)
	if err != nil {
		t.Fatalf("Divide(10, 2) returned unexpected error: %v", err)
	}
	if got != 5 {
		t.Errorf("Divide(10, 2) = %v, want 5", got)
	}

	_, err = Divide(1, 0)
	if err == nil {
		t.Error("Divide(1, 0) should return a non-nil error")
	}
}

func TestMinMax(t *testing.T) {
	min, max := MinMax(3, 1, 4, 1, 5, 9, 2, 6)
	if min != 1 || max != 9 {
		t.Errorf("MinMax(...) = (%d, %d), want (1, 9)", min, max)
	}
}

func TestNewCounter(t *testing.T) {
	counter := NewCounter()
	if counter == nil {
		t.Fatal("NewCounter() returned nil, want a function")
	}
	if got := counter(); got != 1 {
		t.Errorf("first call = %d, want 1", got)
	}
	if got := counter(); got != 2 {
		t.Errorf("second call = %d, want 2", got)
	}
	if got := counter(); got != 3 {
		t.Errorf("third call = %d, want 3", got)
	}

	other := NewCounter()
	if other == nil {
		t.Fatal("second NewCounter() call returned nil, want a function")
	}
	if got := other(); got != 1 {
		t.Errorf("a new counter's first call = %d, want 1 (should be independent)", got)
	}
}
