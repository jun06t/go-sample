package main

import (
	"testing"

	"github.com/jun06t/go-sample/pgv/user"
)

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name    string
		user    *user.User
		wantErr bool
	}{
		{
			name: "有効なユーザー",
			user: &user.User{
				Id:    1,
				Name:  "test_user",
				Email: "test@example.com",
				Age:   25,
				Address: &user.Address{
					Country:    "JP",
					Prefecture: "東京都",
					City:       "渋谷区",
				},
			},
			wantErr: false,
		},
		{
			name: "無効なID（負の数）",
			user: &user.User{
				Id:    -1,
				Name:  "test_user",
				Email: "test@example.com",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "無効な名前（短すぎる）",
			user: &user.User{
				Id:    1,
				Name:  "ab",
				Email: "test@example.com",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "無効なメールアドレス",
			user: &user.User{
				Id:    1,
				Name:  "test_user",
				Email: "invalid-email",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "無効な年齢（未成年）",
			user: &user.User{
				Id:    1,
				Name:  "test_user",
				Email: "test@example.com",
				Age:   17,
			},
			wantErr: true,
		},
		{
			name: "無効な電話番号",
			user: &user.User{
				Id:    1,
				Name:  "test_user",
				Email: "test@example.com",
				Age:   25,
				Phone: "123-456-789",
			},
			wantErr: true,
		},
		{
			name: "有効な電話番号",
			user: &user.User{
				Id:    1,
				Name:  "test_user",
				Email: "test@example.com",
				Age:   25,
				Phone: "0901234567",
				Address: &user.Address{
					Country:    "JP",
					Prefecture: "東京都",
					City:       "渋谷区",
				},
			},
			wantErr: false,
		},
		{
			name: "無効な住所（国コードが長すぎる）",
			user: &user.User{
				Id:    1,
				Name:  "test_user",
				Email: "test@example.com",
				Age:   25,
				Address: &user.Address{
					Country:    "JPN",
					Prefecture: "東京都",
					City:       "渋谷区",
				},
			},
			wantErr: true,
		},
		{
			name: "有効な住所",
			user: &user.User{
				Id:    1,
				Name:  "test_user",
				Email: "test@example.com",
				Age:   25,
				Address: &user.Address{
					Country:    "JP",
					Prefecture: "東京都",
					City:       "渋谷区",
					PostalCode: "1500001",
				},
			},
			wantErr: false,
		},
		{
			name: "無効なタグ（多すぎる）",
			user: &user.User{
				Id:    1,
				Name:  "test_user",
				Email: "test@example.com",
				Age:   25,
				Tags:  []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6"},
			},
			wantErr: true,
		},
		{
			name: "有効なタグ",
			user: &user.User{
				Id:    1,
				Name:  "test_user",
				Email: "test@example.com",
				Age:   25,
				Tags:  []string{"golang", "protobuf", "validation"},
				Address: &user.Address{
					Country:    "JP",
					Prefecture: "東京都",
					City:       "渋谷区",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddressValidation(t *testing.T) {
	tests := []struct {
		name    string
		address *user.Address
		wantErr bool
	}{
		{
			name: "有効な住所",
			address: &user.Address{
				Country:    "JP",
				Prefecture: "東京都",
				City:       "渋谷区",
				PostalCode: "1500001",
			},
			wantErr: false,
		},
		{
			name: "無効な国コード（長すぎる）",
			address: &user.Address{
				Country:    "JPN",
				Prefecture: "東京都",
				City:       "渋谷区",
			},
			wantErr: true,
		},
		{
			name: "無効な国コード（短すぎる）",
			address: &user.Address{
				Country:    "J",
				Prefecture: "東京都",
				City:       "渋谷区",
			},
			wantErr: true,
		},
		{
			name: "無効な郵便番号",
			address: &user.Address{
				Country:    "JP",
				Prefecture: "東京都",
				City:       "渋谷区",
				PostalCode: "123-4567",
			},
			wantErr: true,
		},
		{
			name: "有効な郵便番号",
			address: &user.Address{
				Country:    "JP",
				Prefecture: "東京都",
				City:       "渋谷区",
				PostalCode: "1500001",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.address.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Address.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
