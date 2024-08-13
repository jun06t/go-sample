package mycmp_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	mycmp "github.com/jun06t/go-sample/go-cmp"
	mock_cmp "github.com/jun06t/go-sample/go-cmp/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	type in struct {
		name    string
		age     int
		address *mycmp.Address
	}
	tests := []struct {
		name string
		in   in
		out  mycmp.User
	}{
		{
			name: "success",
			in: in{
				name: "alice",
				age:  20,
				address: &mycmp.Address{
					ZipCode: "102-0101",
					Pref:    "Tokyo",
					City:    "Chiyoda",
					Street:  "Kokyo 1-1-1",
				},
			},
			out: mycmp.User{
				Name: "alice",
				Age:  20,
				Address: mycmp.Address{
					ZipCode: "102-0101",
					Pref:    "Tokyo",
					City:    "Chiyoda",
					Street:  "Kokyo",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := mycmp.NewUser(tt.in.name, tt.in.age, tt.in.address)
			if diff := cmp.Diff(tt.out, out, cmpopts.IgnoreFields(mycmp.User{}, "ID", "Address.Street")); diff != "" {
				t.Errorf("Mismatch (-expected +actual):\n%s", diff)
			}
		})
	}
}

func TestUserService_Save(t *testing.T) {
	type in struct {
		name    string
		age     int
		address *mycmp.Address
	}
	tests := []struct {
		name     string
		injector func(*mock_cmp.MockUserRepository)
		in       in
		err      error
	}{
		{
			name: "success",
			injector: func(m *mock_cmp.MockUserRepository) {
				m.EXPECT().Save(userMatcher{
					expected: mycmp.User{
						Name: "test",
						Age:  20,
						Address: mycmp.Address{
							ZipCode: "102-0101",
							Pref:    "Tokyo",
							City:    "Chiyoda",
							Street:  "Kokyo 1-1-1",
						},
					}}).Return(nil)
			},
			in: in{
				name: "test",
				age:  20,
				address: &mycmp.Address{
					ZipCode: "102-0101",
					Pref:    "Tokyo",
					City:    "Chiyoda",
					Street:  "Kokyo",
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := mock_cmp.NewMockUserRepository(ctrl)
			tt.injector(m)

			s := mycmp.NewUserService(m)
			err := s.Save(tt.in.name, tt.in.age, tt.in.address)
			assert.Equal(t, tt.err, err)
		})
	}
}

type userMatcher struct {
	expected mycmp.User
}

func (m userMatcher) Matches(x interface{}) bool {
	actual, ok := x.(mycmp.User)
	if !ok {
		return false
	}

	// 一部のフィールドを無視して他のフィールドを比較
	return cmp.Equal(m.expected, actual, cmp.FilterPath(func(p cmp.Path) bool {
		if p.String() == "ID" {
			return true
		}
		if p.String() == "Address.Street" {
			return true
		}
		return false
	}, cmp.Ignore()))
}

func (m userMatcher) String() string {
	return "matches User with ID ignored"
}
