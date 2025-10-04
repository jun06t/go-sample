package main

import (
	"fmt"
	"log"

	"github.com/jun06t/go-sample/pgv/user"
	"google.golang.org/protobuf/proto"
)

func main() {
	fmt.Println("=== proto-gen-validate サンプルプログラム ===")
	fmt.Println()

	// 有効なユーザーデータのテスト
	fmt.Println("1. 有効なユーザーデータのテスト:")
	validUser := createValidUser()
	if err := validUser.Validate(); err != nil {
		log.Printf("バリデーションエラー: %v", err)
	} else {
		fmt.Println("✅ バリデーション成功: 有効なユーザーデータ")
		printUser(validUser)
	}
	fmt.Println()

	// 無効なユーザーデータのテスト
	fmt.Println("2. 無効なユーザーデータのテスト:")
	invalidUsers := createInvalidUsers()
	for i, invalidUser := range invalidUsers {
		fmt.Printf("テストケース %d:\n", i+1)
		if err := invalidUser.Validate(); err != nil {
			fmt.Printf("❌ バリデーションエラー: %v\n", err)
		} else {
			fmt.Println("✅ バリデーション成功（予期しない結果）")
		}
		fmt.Println()
	}

	// バリデーション関数のデモ
	fmt.Println("3. バリデーション関数のデモ:")
	demonstrateValidation()
}

// 有効なユーザーデータを作成
func createValidUser() *user.User {
	return &user.User{
		Id:    12345,
		Name:  "john_doe",
		Email: "john.doe@example.com",
		Age:   25,
		Phone: "0901234567",
		Address: &user.Address{
			Country:     "JP",
			Prefecture:  "東京都",
			City:        "渋谷区",
			Street:      "1-2-3 サンプルマンション",
			PostalCode:  "1500001",
		},
		Tags: []string{"golang", "protobuf", "validation"},
	}
}

// 無効なユーザーデータを作成
func createInvalidUsers() []*user.User {
	return []*user.User{
		// 無効なID（負の数）
		{
			Id:    -1,
			Name:  "test_user",
			Email: "test@example.com",
			Age:   25,
		},
		// 無効な名前（短すぎる）
		{
			Id:    1,
			Name:  "ab",
			Email: "test@example.com",
			Age:   25,
		},
		// 無効なメールアドレス
		{
			Id:    1,
			Name:  "test_user",
			Email: "invalid-email",
			Age:   25,
		},
		// 無効な年齢（未成年）
		{
			Id:    1,
			Name:  "test_user",
			Email: "test@example.com",
			Age:   17,
		},
		// 無効な電話番号
		{
			Id:    1,
			Name:  "test_user",
			Email: "test@example.com",
			Age:   25,
			Phone: "123-456-789",
		},
		// 無効な住所（国コードが長すぎる）
		{
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
		// 無効なタグ（多すぎる）
		{
			Id:    1,
			Name:  "test_user",
			Email: "test@example.com",
			Age:   25,
			Tags:  []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6"},
		},
	}
}

// バリデーション関数のデモ
func demonstrateValidation() {
	// ユーザー作成関数
	createUser := func(id int64, name, email string, age int32) *user.User {
		return &user.User{
			Id:    id,
			Name:  name,
			Email: email,
			Age:   age,
		}
	}

	// バリデーション関数
	validateUser := func(u *user.User) error {
		return u.Validate()
	}

	// テストケース
	testCases := []struct {
		name string
		user *user.User
	}{
		{"有効なユーザー", createUser(1, "valid_user", "valid@example.com", 25)},
		{"無効なID", createUser(0, "valid_user", "valid@example.com", 25)},
		{"無効な名前", createUser(1, "ab", "valid@example.com", 25)},
		{"無効なメール", createUser(1, "valid_user", "invalid-email", 25)},
		{"無効な年齢", createUser(1, "valid_user", "valid@example.com", 15)},
	}

	for _, tc := range testCases {
		fmt.Printf("テスト: %s\n", tc.name)
		if err := validateUser(tc.user); err != nil {
			fmt.Printf("  ❌ バリデーションエラー: %v\n", err)
		} else {
			fmt.Printf("  ✅ バリデーション成功\n")
		}
	}
}

// ユーザー情報を表示
func printUser(u *user.User) {
	fmt.Printf("  ID: %d\n", u.Id)
	fmt.Printf("  名前: %s\n", u.Name)
	fmt.Printf("  メール: %s\n", u.Email)
	fmt.Printf("  年齢: %d\n", u.Age)
	if u.Phone != "" {
		fmt.Printf("  電話: %s\n", u.Phone)
	}
	if u.Address != nil {
		fmt.Printf("  住所: %s %s %s\n", u.Address.Country, u.Address.Prefecture, u.Address.City)
	}
	if len(u.Tags) > 0 {
		fmt.Printf("  タグ: %v\n", u.Tags)
	}
