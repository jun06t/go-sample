package foo

import "github.com/google/wire"

func initializeBarBaz() int {
	wire.Build(
		ProvideFoo,
		ProvideBar,
		ProvideBaz,
		NewBarBaz,
	)
	return 0
}
