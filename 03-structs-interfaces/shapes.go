package shapes

import "math"

// TODO: implement each method/function below.
// Run `go test ./03-structs-interfaces/...` to check your work.

// Rectangle is a plain data struct — think of it like a TS
// `type Rectangle = { width: number; height: number }`.
type Rectangle struct {
	Width, Height float64
}

// Area computes the rectangle's area (Width * Height).
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Circle is another plain data struct.
type Circle struct {
	Radius float64
}

// Area computes the circle's area (pi * r^2). Use math.Pi.
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Shape is satisfied by any type with an Area() float64 method — no
// `implements` keyword needed. Rectangle and Circle already satisfy it.
type Shape interface {
	Area() float64
}

// TotalArea returns the sum of Area() across all shapes.
func TotalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}
