# 06-http-basics

Package `httpbasics`. Topic: `net/http` handlers, reading query params,
`encoding/json` responses — the first building block toward a real backend
(Controller layer). Tests use `net/http/httptest` to call handlers directly
without spinning up a real server.

## To implement (`http_basics.go`)

- `HelloHandler(w http.ResponseWriter, r *http.Request)` — read the `name`
  query parameter from the request URL and write `"Hello, <name>!"` as plain
  text. Default to `"World"` if `name` is missing or empty.
- `UserHandler(w http.ResponseWriter, r *http.Request)` — write a JSON
  response for `User{ID: 1, Name: "Alice"}`. Set the `Content-Type` header
  to `application/json` before writing the body.

## Gotchas specific to this batch

- `http.ResponseWriter` is an interface, not something you construct — the
  test hands you one via `httptest.NewRecorder()`. Order matters: any header
  you set (`w.Header().Set(...)`) must happen *before* the first call that
  writes the body (`w.Write(...)` / `json.NewEncoder(w).Encode(...)`) —
  once bytes are written, headers are already locked in and sent.
- `r.URL.Query().Get("name")` returns `""` (empty string) if the param is
  missing — no error, no panic. Same "zero value, not an exception" pattern
  you've already seen with maps.
- Struct tags (`` `json:"id"` ``) are new syntax: they control the JSON key
  name `encoding/json` uses for that field, independent of the Go field name.

```bash
go test ./06-http-basics/...
```

See `notes.md` for the concept primer.
