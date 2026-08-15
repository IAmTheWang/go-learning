# 03-structs-interfaces

Package `shapes`. Topic: struct, method receivers, interfaces — compared to
TS `type`/`interface` and structural typing (Go interfaces are satisfied
implicitly, no `implements` keyword, same spirit as TS structural typing but
Go takes it further — there isn't even a way to declare intent to implement).

## To implement (`shapes.go`)

- `Rectangle{Width, Height float64}.Area() float64` — value receiver
  (`func (r Rectangle) Area() float64`), returns `Width * Height`.
- `Circle{Radius float64}.Area() float64` — value receiver, returns
  `math.Pi * Radius * Radius`.
- `Shape` interface (`Area() float64`) — already defined; `Rectangle` and
  `Circle` satisfy it automatically once their `Area()` methods exist.
- `TotalArea(shapes []Shape) float64` — sum of `Area()` across the slice.

## Gotchas specific to this batch

- These methods use **value receivers** (`(r Rectangle)`, not
  `(r *Rectangle)`) — the receiver is a copy, same semantics as passing the
  struct by value to a regular function. Fine here since `Area()` only reads
  fields; a method that needs to *mutate* the receiver would need a pointer
  receiver instead (see the pointer discussion in the root `LEARNING-LOG.md`).
- No `class ... implements Shape` anywhere — a type satisfies an interface
  purely by having the right method signatures. This is stricter/more
  implicit than TS: there's no explicit annotation you can point to that says
  "Rectangle implements Shape."

See `notes.md` for the full JS/TS comparison.
