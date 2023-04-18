package main

import (
	"net/http"
	"strings"
)

type Route struct {
	Path    string `json:"path"`
	handler http.HandlerFunc
}

type Node[T any] struct {
	Part     string     `json:"part"`
	Children []*Node[T] `json:"children"`
	IsWild   bool       `json:"isWild"`
	Route    T          `json:"route"`
}

func (n *Node[T]) insert(pattern string, route T) {
	parts := strings.Split(pattern, "/")[1:]

	for _, part := range parts {
		child := n.matchChild(part)
		if child == nil {
			child = &Node[T]{
				Part:   part,
				IsWild: part[0] == ':' || part[0] == '*' || part[0] == '{',
			}
			n.Children = append(n.Children, child)
		}
		n = child
	}

	n.Route = route
}

func (n *Node[T]) search(path string) (ret T) {
	parts := strings.Split(path, "/")[1:]

	for _, part := range parts {
		child := n.matchChild(part)
		if child == nil {
			return ret
		}
		n = child
	}

	return n.Route
}

func (n *Node[T]) matchChild(part string) *Node[T] {
	for i := range n.Children {
		if n.Children[i].Part == part || n.Children[i].IsWild {
			return n.Children[i]
		}
	}
	return nil
}
