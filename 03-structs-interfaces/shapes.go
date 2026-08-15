package shapes

// TODO: implement each method/function below.
// Run `go test ./03-structs-interfaces/...` to check your work.

// Rectangle is a plain data struct — think of it like a TS
// `type Rectangle = { width: number; height: number }`.
type Rectangle struct {
	Width, Height float64
}

// Area computes the rectangle's area (Width * Height).
func (r Rectangle) Area() float64 {
	return 0
}

// Circle is another plain data struct.
type Circle struct {
	Radius float64
}

// Area computes the circle's area (pi * r^2). Use math.Pi.
func (c Circle) Area() float64 {
	return 0
}

// Shape is satisfied by any type with an Area() float64 method — no
// `implements` keyword needed. Rectangle and Circle already satisfy it.
type Shape interface {
	Area() float64
}

// TotalArea returns the sum of Area() across all shapes.
func TotalArea(shapes []Shape) float64 {
	return 0
}
