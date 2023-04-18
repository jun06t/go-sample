package main

import (
	"regexp"
	"testing"
)

var (
	patterns = make([]Pattern, 0)
	tree     = &Node[Route]{}
)

func init() {
	routes := []Route{
		{Path: "/hello", handler: Hello},
		{Path: "/hello/{name}", handler: Hello},
		{Path: "/hello/{name}/foo", handler: Hello},
		{Path: "/foo", handler: Hello},
		{Path: "/foo/bar", handler: Hello},
		{Path: "/foo/bar/{name}", handler: Hello},
		{Path: "/bar", handler: Hello},
		{Path: "/bar/baz", handler: Hello},
		{Path: "/bar/baz/{name}", handler: Hello},
		{Path: "/bar/baz/{name}/hoge", handler: Hello},
	}
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
}

func BenchmarkRegexFirst(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := range patterns {
			if ok := patterns[j].Pattern.MatchString("/hello"); ok {
				continue
			}
		}
	}
}

func BenchmarkRadixFirst(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tree.search("/hello")
	}
}

func BenchmarkRegexMid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := range patterns {
			if ok := patterns[j].Pattern.MatchString("/foo/bar/12345"); ok {
				continue
			}
		}
	}
}

func BenchmarkRadixMid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tree.search("/foo/bar/12345")
	}
}

func BenchmarkRegexLast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := range patterns {
			if ok := patterns[j].Pattern.MatchString("/bar/baz/12345/hoge"); ok {
				continue
			}
		}
	}
}

func BenchmarkRadixLast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tree.search("/bar/baz/12345/hoge")
	}
}

func BenchmarkRegexRound(b *testing.B) {
	path := genPath()
	for i := 0; i < b.N; i++ {
		for j := range patterns {
			if ok := patterns[j].Pattern.MatchString(path()); ok {
				continue
			}
		}
	}
}

func BenchmarkRadixRound(b *testing.B) {
	path := genPath()
	for i := 0; i < b.N; i++ {
		tree.search(path())
	}
}

func genPath() func() string {
	paths := []string{
		"/hello",
		"/hello/12345",
		"/hello/12345/foo",
		"/foo",
		"/foo/bar",
		"/foo/bar/12345",
		"/bar",
		"/bar/baz",
		"/bar/baz/12345",
		"/bar/baz/12345/hoge",
	}
	count := 0
	max := len(paths)

	return func() string {
		if count >= max {
			count = 0
		}
		path := paths[count]
		count++
		return path
	}
}
