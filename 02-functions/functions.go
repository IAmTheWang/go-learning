package functions

import "errors"

// TODO: implement each function below.
// Run `go test ./02-functions/...` to check your work.

// Divide returns a/b. If b is 0, it must return an error instead of
// panicking or returning Inf/NaN. Hint: errors.New("division by zero").
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// MinMax returns the smallest and largest values among nums.
// nums is variadic: call it as MinMax(1, 2, 3) or MinMax(slice...).
// min and max are named return values — a bare `return` sends back
// whatever they currently hold.
func MinMax(nums ...int) (min, max int) {
	a, b := nums[0], nums[0]
	for index := range nums[1:] {

		if nums[index] > b {
			b = nums[index]
		}
		if nums[index] < a {
			a = nums[index]

		}

	}
	return a, b

}

// NewCounter returns a function that, each time it's called, returns an
// incrementing count starting at 1. Two counters returned by separate
// calls to NewCounter must be independent of each other.
func NewCounter() func() int {
	return nil
}
