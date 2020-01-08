package foo

import "github.com/google/wire"

func initializeFooBarBaz() *FooBarBaz {
	wire.Build(
		ProvideFoo,
		ProvideBar,
		wire.Struct(new(FooBarBaz), "MyFoo", "MyBar"),
	)
	return &FooBarBaz{}
}
