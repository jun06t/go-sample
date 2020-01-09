package foo

import (
	"github.com/google/wire"
	"github.com/jun06t/go-sample/wire/packages/model"
)

func initializeFooBarBaz() int {
	wire.Build(model.Set)
	return 0
}
