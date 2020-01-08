package foo

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

type FooBarBaz struct {
	MyFoo Foo
	MyBar Bar
	MyBaz Baz
}
