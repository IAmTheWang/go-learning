# 07-errors-interfaces

Package `errorsinterfaces`. Topic: custom error types, error wrapping
(`fmt.Errorf` + `%w`), `errors.As`, and the `fmt.Stringer` interface —
idiomatic Go error handling beyond the basic `if err != nil` pattern already
seen in `01-basics`/`06-http-basics`.

## To implement (`errors_interfaces.go`)

- `(*NotFoundError).Error() string` — return a message like
  `"item 42 not found"` using the struct's `ID` field.
- `FindItem(items map[int]string, id int) (string, error)` — look up `id`;
  return `&NotFoundError{ID: id}` if missing.
- `LoadItemConfig(items map[int]string, id int) error` — call `FindItem`,
  and on error wrap it with `fmt.Errorf("load config: %w", err)`.
- `IsNotFound(err error) bool` — use `errors.As` with a `*NotFoundError`
  target to check whether `err` is or wraps one, even through the
  `LoadItemConfig` wrapping layer.
- `(Temperature).String() string` — format as `"23.5°C"` (one decimal place)
  to satisfy `fmt.Stringer`.

## Gotchas specific to this batch

- Stub bodies intentionally import nothing beyond what's already there —
  add `"fmt"` and `"errors"` to the import block yourself as you implement
  each piece (see the root `CLAUDE.md`/`LEARNING-LOG.md` note on why stubs
  don't pre-import packages they don't use yet: unused imports are a Go
  compile error, so a stub that already imported `errors` before you touch
  `IsNotFound` would fail to build).
- Plain type assertion (`err.(*NotFoundError)`) breaks once an error has
  been wrapped by `%w` — you must use `errors.As`, which walks the
  `Unwrap()` chain instead of checking the error's immediate concrete type.
- `errors.As` takes a pointer to your target variable (`&nf`, where
  `nf` is `*NotFoundError`) so it can write the match back to it — same
  "pass a pointer so the callee can mutate the caller's variable" pattern
  as `05-goroutines-channels`, applied to a different problem.
- `Temperature` uses a **value receiver** (`(t Temperature)`), matching
  `03-structs-interfaces`'s `Rectangle`/`Circle` — `String()` only reads the
  value, no mutation needed.

```bash
go test ./07-errors-interfaces/...
```

See `notes.md` for the concept primer.
