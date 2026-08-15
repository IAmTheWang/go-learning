# 02-functions

Package `functions`. Topic: multiple return values, named return values,
variadic parameters, closures — compared to JS destructured returns,
`...args`, and closures (which work the same way in both languages).

## To implement (`functions.go`)

- `Divide(a, b float64) (float64, error)` — must return an error via
  `errors.New("division by zero")` when `b == 0`, not panic and not return
  `Inf`/`NaN`. This is Go's `(value, error)` idiom in place of try/catch.
- `MinMax(nums ...int) (min, max int)` — variadic input, **named** return
  values (`min`, `max` are already declared by the signature; a bare
  `return` sends back whatever they currently hold).
- `NewCounter() func() int` — returns a closure; each call to the returned
  function increments and returns a count starting at 1. Two counters from
  separate `NewCounter()` calls must be independent (separate captured state).

## Gotcha: `functions_test.go`'s nil-guard is intentional, don't remove it

`TestNewCounter` checks `if counter == nil { t.Fatal(...) }` before calling
`counter()`. This exists because the initial TODO stub returns `nil`, and
calling a `nil` function value directly panics with a scary segfault-style
trace instead of a clean test failure. Once `NewCounter` is implemented
correctly the guard is a no-op — it's defensive, not a bug.

See `notes.md` for the JS/TS comparison, and the repo root's `LEARNING-LOG.md`
for deeper explanations (pointers, zero values, etc.) that came up while
working through this repo.
