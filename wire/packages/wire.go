package foo

import "github.com/google/wire"

func initializeFooBarBaz() int {
	wire.Build(Set)
	return 0
}
