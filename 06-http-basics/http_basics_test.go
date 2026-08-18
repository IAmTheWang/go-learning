package httpbasics

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want string
	}{
		{"with name", "/hello?name=Alice", "Hello, Alice!"},
		{"missing name", "/hello", "Hello, World!"},
	}
	for _, c := range cases {
		req := httptest.NewRequest("GET", c.url, nil)
		rec := httptest.NewRecorder()

		HelloHandler(rec, req)

		if got := rec.Body.String(); got != c.want {
			t.Errorf("%s: HelloHandler(%s) body = %q, want %q", c.name, c.url, got, c.want)
		}
	}
}

func TestUserHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/user", nil)
	rec := httptest.NewRecorder()

	UserHandler(rec, req)

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var got User
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	want := User{ID: 1, Name: "Alice"}
	if got != want {
		t.Errorf("UserHandler body = %+v, want %+v", got, want)
	}
}
