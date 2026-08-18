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

## Real mistakes made while implementing this batch (keep an eye out)

- **`if init; condition { }` is not "condition; then run this statement."**
  The part before `;` is an ordinary statement that runs first (usually a
  short variable declaration); only the part after `;` is the actual boolean
  condition, and the body always needs `{ }` regardless. Both `Divide` and
  `MinMax` first-draft attempts wrote something shaped like
  `if judge := b == 0; return ...` (treating `;` as "and then do this"),
  which doesn't compile.
- **`MinMax`'s min/max trackers must start at `nums[0]`, not `0`.** Starting
  both accumulators at the zero value looks reasonable but silently breaks
  the moment every element is `>= 0` (min never updates away from 0) or every
  element is `<= 0` (max never updates away from 0) — the test input
  `MinMax(3, 1, 4, 1, 5, 9, 2, 6)` happens to be all-positive, so this bug
  reproduces on the very first test case, not some obscure edge case.
- **No implicit `int`↔`float64` conversion, even when it "should obviously
  work."** Declaring `var a, b float64` to hold running min/max while
  `nums` is `[]int` fails to compile
  (`invalid operation: nums[index] > b (mismatched types int and float64)`)
  — the accumulator type must match the slice's element type exactly.

See `notes.md` for the JS/TS comparison, and the repo root's `LEARNING-LOG.md`
for deeper explanations (pointers, zero values, etc.) that came up while
working through this repo.
