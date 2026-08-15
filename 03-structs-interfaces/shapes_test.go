package shapes

import "testing"

func approxEqual(a, b float64) bool {
	const epsilon = 0.001
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

func TestRectangleArea(t *testing.T) {
	r := Rectangle{Width: 3, Height: 4}
	if got := r.Area(); got != 12 {
		t.Errorf("Rectangle{3,4}.Area() = %v, want 12", got)
	}
}

func TestCircleArea(t *testing.T) {
	c := Circle{Radius: 2}
	got := c.Area()
	want := 12.566370614359172
	if !approxEqual(got, want) {
		t.Errorf("Circle{2}.Area() = %v, want ~%v", got, want)
	}
}

func TestTotalArea(t *testing.T) {
	shapes := []Shape{
		Rectangle{Width: 2, Height: 3},
		Circle{Radius: 1},
	}
	got := TotalArea(shapes)
	want := 6 + 3.14159265
	if !approxEqual(got, want) {
		t.Errorf("TotalArea(...) = %v, want ~%v", got, want)
	}
}
