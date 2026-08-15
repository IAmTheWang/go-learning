package concurrency

// TODO: implement each function/type below.
// Run `go test ./05-goroutines-channels/...` to check your work, then
// `go test -race ./05-goroutines-channels/...` to confirm SafeCounter
// has no data races.

// SquareAsync spawns a goroutine that computes n*n and sends the result
// on the returned channel. The caller receives from the channel to get
// the result (see the notes.md example using `select` + `time.After`
// for how the test protects itself from a hang).
func SquareAsync(n int) <-chan int {
	return nil
}

// SumConcurrent splits nums into `workers` roughly-equal chunks, sums each
// chunk in its own goroutine, and returns the total. Use sync.WaitGroup to
// wait for all goroutines to finish, and a channel to collect the partial
// sums safely. Must handle workers <= 0 or workers > len(nums).
func SumConcurrent(nums []int, workers int) int {
	return 0
}

// SafeCounter can be incremented from many goroutines at once without a
// data race. Protect the internal count with a sync.Mutex.
type SafeCounter struct {
	// TODO: add fields (a sync.Mutex and an int)
}

// Increment adds 1 to the counter. Safe to call from many goroutines
// concurrently.
func (c *SafeCounter) Increment() {
	// TODO
}

// Value returns the current count. Safe to call while other goroutines
// are calling Increment.
func (c *SafeCounter) Value() int {
	return 0
}
