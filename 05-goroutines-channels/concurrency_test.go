package concurrency

import (
	"sync"
	"testing"
	"time"
)

func TestSquareAsync(t *testing.T) {
	ch := SquareAsync(5)
	select {
	case got := <-ch:
		if got != 25 {
			t.Errorf("SquareAsync(5) = %d, want 25", got)
		}
	case <-time.After(time.Second):
		t.Fatal("SquareAsync(5) timed out waiting for a result on the channel")
	}
}

func TestSumConcurrent(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	want := 55
	for _, workers := range []int{1, 2, 3, 4, 10, 100} {
		if got := SumConcurrent(nums, workers); got != want {
			t.Errorf("SumConcurrent(nums, %d workers) = %d, want %d", workers, got, want)
		}
	}
}

func TestSumConcurrentEmpty(t *testing.T) {
	if got := SumConcurrent(nil, 4); got != 0 {
		t.Errorf("SumConcurrent(nil, 4) = %d, want 0", got)
	}
}

func TestSafeCounter(t *testing.T) {
	c := &SafeCounter{}
	const n = 1000

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}
	wg.Wait()

	if got := c.Value(); got != n {
		t.Errorf("Value() = %d, want %d", got, n)
	}
}
