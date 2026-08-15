package basics

import "testing"

func TestSum(t *testing.T) {
	cases := []struct {
		a, b, want int
	}{
		{1, 2, 3},
		{-1, 1, 0},
		{0, 0, 0},
	}
	for _, c := range cases {
		if got := Sum(c.a, c.b); got != c.want {
			t.Errorf("Sum(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestIsEven(t *testing.T) {
	cases := []struct {
		n    int
		want bool
	}{
		{2, true},
		{3, false},
		{0, true},
		{-4, true},
	}
	for _, c := range cases {
		if got := IsEven(c.n); got != c.want {
			t.Errorf("IsEven(%d) = %v, want %v", c.n, got, c.want)
		}
	}
}

func TestFizzBuzz(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{
		{1, "1"},
		{3, "Fizz"},
		{5, "Buzz"},
		{15, "FizzBuzz"},
		{7, "7"},
	}
	for _, c := range cases {
		if got := FizzBuzz(c.n); got != c.want {
			t.Errorf("FizzBuzz(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

func TestGrade(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{95, "A"},
		{82, "B"},
		{71, "C"},
		{60, "D"},
		{40, "F"},
	}
	for _, c := range cases {
		if got := Grade(c.score); got != c.want {
			t.Errorf("Grade(%d) = %q, want %q", c.score, got, c.want)
		}
	}
}
