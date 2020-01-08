package foo

import "github.com/google/wire"

func initializeFooBarBaz() *FooBarBaz {
	wire.Build(
		ProvideFoo,
		ProvideBar,
		ProvideBaz,
		wire.Struct(new(FooBarBaz), "*"),
	)
	return &FooBarBaz{}
}
