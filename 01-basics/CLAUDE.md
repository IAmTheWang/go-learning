# 01-basics

Package `basics`. Topic: variables, static typing, consts, control flow —
compared to JS/TS `let`/`const`, `if`/`switch`, and implicit type coercion
(Go has none).

## To implement (`basics.go`)

- `Sum(a, b int) int`
- `IsEven(n int) bool`
- `FizzBuzz(n int) string` — divisible by 3 → `"Fizz"`, by 5 → `"Buzz"`, by
  both → `"FizzBuzz"`, else the number itself as a string
  (`strconv.Itoa(n)`).
- `Grade(score int) string` — `>=90` "A", `>=80` "B", `>=70` "C", `>=60` "D",
  else "F". Meant to be written with a `switch` — Go's `switch` does **not**
  fall through by default (opposite of JS/TS), so no `break` needed per case.

## Gotchas specific to this batch

- No implicit int↔float or int↔string conversion — `strconv.Itoa`/`strconv.Atoi`
  are required, unlike JS's automatic coercion.
- `switch` without a condition (`switch { case x > 90: ... }`) is idiomatic Go
  for this kind of range check — cleaner than an `if`/`else if` chain.

See `notes.md` for the full JS/TS comparison, and the repo root's
`LEARNING-LOG.md` for deeper "why" explanations that came up in Q&A.
