package foo

type Foo int
type Bar int
type Baz int

func ProvideFoo() Foo {
	return Foo(1)
}

func ProvideBar(f Foo) Bar {
	return Bar(f)
}

func ProvideBaz(f Foo) Baz {
	return Baz(f)
}

func NewBarBaz(bar Bar, baz Baz) int {
	return int(bar) + int(baz)
}
