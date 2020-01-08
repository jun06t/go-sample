// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package foo

// Injectors from wire.go:

func initializeFooBarBaz() *FooBarBaz {
	foo := ProvideFoo()
	bar := ProvideBar()
	fooBarBaz := &FooBarBaz{
		MyFoo: foo,
		MyBar: bar,
	}
	return fooBarBaz
}
