// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package foo

// Injectors from wire.go:

func initializeFooBarBaz() int {
	foo := ProvideFoo()
	bar := ProvideBar()
	baz := ProvideBaz()
	int2 := NewFooBarBaz(foo, bar, baz)
	return int2
}