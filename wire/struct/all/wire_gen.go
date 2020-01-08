// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package foo

// Injectors from wire.go:

func initializeFooBarBaz() *FooBarBaz {
	foo := ProvideFoo()
	bar := ProvideBar()
	baz := ProvideBaz()
	fooBarBaz := &FooBarBaz{
		MyFoo: foo,
		MyBar: bar,
		MyBaz: baz,
	}
	return fooBarBaz
}
