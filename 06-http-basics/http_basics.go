package httpbasics

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TODO: implement each handler below.
// Run `go test ./06-http-basics/...` to check your work.

// HelloHandler writes "Hello, <name>!" as plain text, where <name> comes
// from the "name" query parameter (e.g. /hello?name=Alice).
// If "name" is missing or empty, default to "World".
//
//	app.get('/hello', (req, res) => {
//	  const name = req.query.name || 'World';
//	  res.send(`Hello, ${name}!`);
//	});
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s!", name)))
}

// User is the JSON shape returned by UserHandler.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserHandler writes a JSON response representing a User with ID 1 and
// Name "Alice". Set the "Content-Type" header to "application/json".
// res.setHeader('Content-Type', 'application/json');
// res.send(JSON.stringify({ id: 1, name: 'Alice' }));
// 或者更常见直接用 res.json({ id: 1, name: 'Alice' })，两步合一
func UserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.Header().Set("Content-Type", "application/json")
	user := User{ID: 1, Name: "Alice"}
	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
