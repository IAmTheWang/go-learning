# 04-slices-maps

Package `collections`. Topic: slice, map — compared to JS `Array` and
`Map`/`Object`. See the root `LEARNING-LOG.md` for the deep dive on how Go
slices differ from JS arrays (re-slicing shares the underlying array, unlike
`Array.prototype.slice()` which copies; `append` may or may not write into
shared backing storage depending on capacity).

## To implement (`collections.go`)

- `Reverse(s []int) []int` — return a **new** slice in reverse order.
  Do not mutate `s` itself (careful: because slices can alias their backing
  array, an in-place reverse-then-return would violate this contract if `s`
  and the caller's original slice share storage).
- `Dedup(s []int) []int` — remove duplicates, preserve first-occurrence
  order. Use `map[int]bool` as a set — Go has no built-in `Set` type like
  JS's `Set`.
- `WordCount(s string) map[string]int` — split on whitespace
  (`strings.Fields(s)`), count occurrences per word. Requires adding
  `"strings"` to the import block yourself (deliberately not pre-added).
  Remember: a `nil` map can be read but panics on write — `make(map[string]int)`
  first.

## Gotchas specific to this batch

- Slice aliasing is the single biggest JS→Go surprise in this batch. If
  `Reverse`/`Dedup` build their result by slicing into the input rather than
  building a genuinely new backing array (e.g. via `make` + explicit copy,
  or `append` to a `nil` slice), you can end up silently mutating the
  caller's data.

See `notes.md` for the full JS/TS comparison.
