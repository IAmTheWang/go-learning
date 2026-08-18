# playground

Throwaway `package main` experiments — not an exercise, not graded, no tests.
This exists specifically so one-off "let me just run this and see" code never
gets mixed into the numbered exercise directories (which would break their
build — see the root `CLAUDE.md`'s "one package per directory" rule).

## Rules

- Each demo gets its **own subdirectory** under `playground/`, e.g.
  `playground/foo_demo/foo_demo.go` — every subdirectory is its own
  `package main` with its own `func main()`.
- Run with `go run ./playground/<name>`.
- **Don't add a second top-level `.go` file directly under `playground/`
  alongside an existing one.** Two `package main` files with two `func main()`
  in the *same* directory don't compile — `go build ./...`/`go test ./...`
  recurse into every directory as one package, so this breaks the whole
  repo's test run, not just playground (this happened once: adding
  `layers_demo.go` next to `select_demo.go` broke `go test ./...` with
  `main redeclared in this block` until both were split into their own
  subdirectories).
- Fine to overwrite, add print statements, or delete demos here freely —
  nothing depends on this directory's contents staying stable.

## Current contents

- `select_demo/select_demo.go` — demonstrates `select` with a `default` case
  (non-blocking channel check). It has an inherent goroutine-scheduling race:
  whether the `default` branch or the "received" branch fires depends on
  whether the sender goroutine gets scheduled before the `select` runs. That
  race is the point of the demo, not a bug — see `LEARNING-LOG.md` at the
  repo root for the full explanation.
- `layers_demo/layers_demo.go` — Controller → Service → DAO/Repository
  layered example (get-user-by-id flow), with a Domain struct vs. a DTO
  struct to show why the DTO shape deliberately excludes internal-only
  fields (`PasswordHash`) when crossing the layer boundary. Print statements
  are labeled with the restaurant-staff analogy (服务员/大厨/仓库) discussed
  in chat, so the console output traces 1:1 onto that mental model.
