package foo

import "github.com/google/wire"

func initializeFooBarBaz() int {
	wire.Build(
		NewFooBarBaz,
		ProvideFoo,
		ProvideBar,
		ProvideBaz,
	)
	return 0
}
