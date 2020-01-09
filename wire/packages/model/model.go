package model

import "github.com/google/wire"

type Foo int
type Bar int
type Baz int

func ProvideFoo() Foo {
	return Foo(1)
}

func ProvideBar() Bar {
	return Bar(2)
}

func ProvideBaz() Baz {
	return Baz(3)
}

func NewFooBarBaz(foo Foo, bar Bar, baz Baz) int {
	return int(foo) + int(bar) + int(baz)
}

var Set = wire.NewSet(
	NewFooBarBaz,
	ProvideFoo,
	ProvideBar,
	ProvideBaz,
)
