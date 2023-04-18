package main

import (
	"fmt"
	"net/http"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/hello/")
	fmt.Fprintf(w, "Hello, %s!\n", name)
}

func Hello2(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/hello/")
	fmt.Fprintf(w, "Hello2, %s!\n", name)
}

func main() {
	routes := []Route{
		{Path: "/hello", handler: Index},
		{Path: "/hello/:name", handler: Hello},
		{Path: "/hello/:name/foo", handler: Hello2},
		{Path: "/foo", handler: Index},
	}

	tree := &Node[Route]{}
	for _, route := range routes {
		tree.insert(route.Path, route)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		route := tree.search(path)

		route.handler(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
