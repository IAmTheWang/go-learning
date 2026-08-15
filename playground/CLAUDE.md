# playground

Throwaway `package main` experiments — not an exercise, not graded, no tests.
This exists specifically so one-off "let me just run this and see" code never
gets mixed into the numbered exercise directories (which would break their
build — see the root `CLAUDE.md`'s "one package per directory" rule).

## Rules

- Every file here is `package main` with its own `func main()`.
- Run with `go run playground/<file>.go`.
- Fine to overwrite, add print statements, or delete files here freely —
  nothing depends on this directory's contents staying stable.

## Current contents

- `select_demo.go` — demonstrates `select` with a `default` case (non-blocking
  channel check). It has an inherent goroutine-scheduling race: whether the
  `default` branch or the "received" branch fires depends on whether the
  sender goroutine gets scheduled before the `select` runs. That race is the
  point of the demo, not a bug — see `LEARNING-LOG.md` at the repo root for
  the full explanation.
