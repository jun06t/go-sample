//+build wireinject

package foo

import "github.com/google/wire"

func initializeFooBarBaz() *FooBarBaz {
	wire.Build(
		Set,
		//		ProvideFoo,
		//		ProvideBar,
		//		ProvideBaz,
		//		wire.Struct(new(FooBarBaz), "MyFoo", "MyBar"),
	)
	return &FooBarBaz{}
}

var Set = wire.NewSet(
	ProvideFoo,
	ProvideBar,
	ProvideBaz,
	wire.Struct(new(FooBarBaz), "MyFoo", "MyBar"),
)
