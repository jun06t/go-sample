// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package user

// Injectors from wire.go:

func initializeUserService() *UserService {
	user := NewUserRepository()
	userService := NewUserService(user)
	return userService
}
