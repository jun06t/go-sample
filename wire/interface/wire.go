package user

import "github.com/google/wire"

func initializeUserService() *UserService {
	wire.Build(
		NewUserService,
		NewUserRepository,
		wire.Bind(new(UserRepository), new(*User)),
	)
	return &UserService{}
}
