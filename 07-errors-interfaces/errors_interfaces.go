package errorsinterfaces

// TODO: implement each function/method below.
// Run `go test ./07-errors-interfaces/...` to check your work.

// NotFoundError is a custom error type carrying the missing item's ID.
// Compare to a JS custom error class:
//
//	class NotFoundError extends Error {
//	  constructor(id) {
//	    super(`item ${id} not found`);
//	    this.id = id;
//	  }
//	}
type NotFoundError struct {
	ID int
}

// Error implements the built-in `error` interface (any type with this exact
// method signature IS an error — no "implements" keyword, no inheritance).
func (e *NotFoundError) Error() string {
	// TODO: return a message like "item 42 not found" (use fmt.Sprintf)
	return ""
}

// FindItem looks up id in items. If missing, return a *NotFoundError
// (not a generic errors.New) so callers can later distinguish "not found"
// from other kinds of failure.
//
//	function findItem(items, id) {
//	  if (!(id in items)) throw new NotFoundError(id);
//	  return items[id];
//	}
func FindItem(items map[int]string, id int) (string, error) {
	// TODO
	return "", nil
}

// LoadItemConfig calls FindItem and, on failure, wraps the error with added
// context using fmt.Errorf's %w verb — this preserves the original error so
// it can still be unwrapped later, unlike just formatting a new string.
//
//	// JS equivalent relies on `Error#cause`:
//	try { findItem(items, id) } catch (err) {
//	  throw new Error(`load config`, { cause: err });
//	}
func LoadItemConfig(items map[int]string, id int) error {
	// TODO: call FindItem, and if it errors, return
	// fmt.Errorf("load config: %w", err) — otherwise return nil
	return nil
}

// IsNotFound reports whether err is (or wraps) a *NotFoundError, using
// errors.As to walk the wrapped-error chain instead of a plain type
// assertion (which would fail once the error has been wrapped by %w).
func IsNotFound(err error) bool {
	// TODO: declare `var target *NotFoundError` and use errors.As(err, &target)
	return false
}

// Temperature is a float64 with a custom String() method, satisfying the
// fmt.Stringer interface. Once a type has String() string, fmt.Println and
// %v automatically call it instead of printing the raw float.
//
//	// JS equivalent: overriding toString()
//	class Temperature {
//	  toString() { return `${this.value}°C`; }
//	}
type Temperature float64

// String implements fmt.Stringer.
func (t Temperature) String() string {
	// TODO: return e.g. "23.5°C" — use fmt.Sprintf("%.1f°C", float64(t))
	return ""
}
