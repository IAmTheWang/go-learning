# go-learning

## What this repo is

A hands-on Go curriculum for a frontend developer (TS/React/Vue background,
7 years experience) learning Go by writing code, not by watching tutorials.
Every concept is anchored to a JS/TS equivalent the user already knows.

## Structure

Each numbered directory is a self-contained Go package and one exercise batch:

- `xxx.go` — the non-test file, has `// TODO` stubs to implement.
- `xxx_test.go` — table-driven tests. **Never edit test files to make a TODO
  pass** — fix the implementation instead. Tests are written idiomatically and
  double as a reference for Go's test style.
- `notes.md` — concept primer for that batch, written as a JS/TS comparison.

Directories, in the order they should be done:

| Dir | Topic |
|---|---|
| `01-basics` | variables, types, consts, control flow |
| `02-functions` | multi-return, named returns, variadic, closures |
| `03-structs-interfaces` | struct, method receivers, interfaces (structural typing) |
| `04-slices-maps` | slice, map |
| `05-goroutines-channels` | goroutine, channel, select, WaitGroup, Mutex |
| `06-http-basics` | net/http handlers, query params, encoding/json responses |
| `07-errors-interfaces` | custom error types, error wrapping (`%w`/`errors.As`), `fmt.Stringer` |

Each subdirectory has its own `CLAUDE.md` with exercise-specific notes and gotchas.

## Non-exercise directories

- `playground/` — throwaway `package main` experiments, not graded, no tests.
  See `playground/CLAUDE.md`.
- `LEARNING-LOG.md` — accumulated Q&A notes from working sessions: the "why does
  this even work this way" explanations that don't belong in any single
  exercise's `notes.md`. Check it before re-explaining a concept from scratch;
  append genuinely new insights here rather than scattering them across
  per-directory files.
- `TECH_LEARNING_ROADMAP.md` — how this repo fits into the user's broader job-hunt
  learning plan (Go is a nice-to-have, lower priority than Docker/AWS SAA).

## Critical rule: one package per directory

All `.go` files in the same directory must declare the same `package` name.
Never create a `package main` file inside `01-basics`/`02-functions`/etc. — it
will break the whole directory's build (`found packages X and main in ...`).
One-off runnable experiments belong in `playground/` instead.

## Common commands

```bash
go test ./...                                  # run everything
go test ./01-basics/...                        # run one batch
gofmt -l .                                     # formatting check
go vet ./...                                   # static checks
go test -race ./05-goroutines-channels/...     # race detector (concurrency batch)
```
