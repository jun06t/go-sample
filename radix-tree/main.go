package main

import (
	"fmt"
	"net/http"
	"regexp"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!\n", r.URL.Path)
}

type Route struct {
	Path    string `json:"path"`
	handler http.HandlerFunc
}

type Pattern struct {
	Pattern *regexp.Regexp
	Route   Route
}

type Patterns []Pattern

func (p Patterns) Match(path string) Route {
	for i := range p {
		if ok := p[i].Pattern.MatchString(path); ok {
			return p[i].Route
		}
	}
	return Route{}
}

var reg = regexp.MustCompile("{[a-zA-Z0-9_]*}")

func main() {
	routes := []Route{
		{Path: "/hello", handler: Hello},
		{Path: "/hello/{name}", handler: Hello},
		{Path: "/hello/{name}/foo", handler: Hello},
	}

	tree := &Node[Route]{}
	patterns := make(Patterns, 0, len(routes))
	for _, route := range routes {
		// Radix tree
		tree.insert(route.Path, route)

		// Regexp
		pt := "^" + reg.ReplaceAllString(route.Path, "[-a-zA-Z0-9_]*") + "$"
		patterns = append(patterns, Pattern{
			Pattern: regexp.MustCompile(pt),
			Route:   route,
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// radix tree search
		route := tree.search(path)

		// regexp loop
		//route := patterns.Match(path)
		route.handler(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
