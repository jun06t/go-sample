//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
package mycmp

import "github.com/google/uuid"

type User struct {
	ID      string
	Name    string
	Age     int
	Address Address
}
type Address struct {
	ZipCode string
	Pref    string
	City    string
	Street  string
}

func NewUser(name string, age int, addr *Address) User {
	u := User{
		ID:   uuid.New().String(),
		Name: name,
		Age:  age,
	}
	if addr != nil {
		u.Address = *addr
	}
	return u
}

type UserRepository interface {
	Save(User) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return UserService{repo: repo}
}

func (s UserService) Save(name string, age int, addr *Address) error {
	// UUIDを生成してIDにセット
	u := NewUser(name, age, addr)

	// ユーザー情報を保存
	err := s.repo.Save(u)
	if err != nil {
		return err
	}

	return nil
}
