package collections

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	original := []int{1, 2, 3}
	got := Reverse(original)
	want := []int{3, 2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Reverse(%v) = %v, want %v", original, got, want)
	}
	if !reflect.DeepEqual(original, []int{1, 2, 3}) {
		t.Errorf("Reverse must not modify its input, got %v", original)
	}
}

func TestDedup(t *testing.T) {
	got := Dedup([]int{1, 2, 2, 3, 1, 4})
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Dedup(...) = %v, want %v", got, want)
	}
}

func TestWordCount(t *testing.T) {
	got := WordCount("the quick brown fox the lazy fox")
	want := map[string]int{
		"the":   2,
		"quick": 1,
		"brown": 1,
		"fox":   2,
		"lazy":  1,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordCount(...) = %v, want %v", got, want)
	}
}
