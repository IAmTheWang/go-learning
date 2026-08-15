# 05-goroutines-channels

Package `concurrency`. Topic: goroutines, channels, `select`,
`sync.WaitGroup`, `sync.Mutex` — compared to JS's single-threaded
`async`/`await`. This batch has by far the most Q&A behind it; see the root
`LEARNING-LOG.md` for the full treatment of blocking vs non-blocking,
`select`, channel mechanics, and why `make()`'s signature looks the way it
does.

## To implement (`concurrency.go`)

- `SquareAsync(n int) <-chan int` — spawn a goroutine that computes `n*n` and
  sends it on the returned channel.
- `SumConcurrent(nums []int, workers int) int` — split `nums` into `workers`
  chunks, sum each chunk in its own goroutine, combine results. Must handle
  `workers <= 0` or `workers > len(nums)`. Use `sync.WaitGroup` to wait for
  all goroutines and a channel (or a mutex-protected accumulator) to collect
  partial sums safely.
- `SafeCounter` — add a `sync.Mutex` and an `int` field; `Increment()` and
  `Value()` must be safe to call from many goroutines concurrently.

## Gotchas specific to this batch

- `TestSquareAsync` wraps its channel receive in `select` + `case
  <-time.After(time.Second)` specifically so an unimplemented (`nil`-channel)
  stub fails fast with a clean timeout message instead of hanging the test
  suite forever. Don't remove this — it's a safety net, not incidental.
- Verify with the race detector, not just `go test`:

  ```bash
  go test -race ./05-goroutines-channels/...
  ```

  `SafeCounter` without a real mutex will pass a plain `go test` run
  sometimes and still be broken — `-race` is what actually catches it.

See `notes.md` for the concept primer.
