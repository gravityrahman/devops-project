package main

import (
	"fmt"
	"net/http"
)

// helloHandler Prints "Hello, world!" to the response writer on GET requests.
// On all non GET requests returns a http.StatusMethodNotAllowed.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		_, _ = fmt.Fprint(w, "Hello, World!")
	default:
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}
