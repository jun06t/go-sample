# proto-gen-validate サンプルプロジェクト

このプロジェクトは、[proto-gen-validate](https://github.com/envoyproxy/protoc-gen-validate)を使用したGo言語のサンプルコードです。Protocol Buffersのメッセージに対してバリデーションルールを定義し、自動的にバリデーションコードを生成する方法をデモンストレーションします。

## 概要

proto-gen-validateは、Protocol Buffersのメッセージ定義にバリデーションルールを追加し、Go言語のバリデーションコードを自動生成するツールです。これにより、型安全で効率的なデータバリデーションを実現できます。

## プロジェクト構成

```
.
├── user.proto          # Protocol Buffersの定義ファイル（バリデーションルール付き）
├── go.mod              # Goモジュール定義
├── Makefile            # ビルド・実行用のMakefile
├── main.go             # メインサンプルプログラム
├── user_test.go        # テストコード
└── README.md           # このファイル
```

## 機能

### 定義されたバリデーションルール

#### User メッセージ
- **ID**: 正の整数（必須）
- **名前**: 3-20文字の英数字とアンダースコア（必須）
- **メールアドレス**: 有効なメール形式（必須）
- **年齢**: 18-120歳（必須）
- **電話番号**: 10桁の数字（オプション）
- **住所**: Addressメッセージ（オプション）
- **タグ**: 最大5個、各1-10文字（オプション）

#### Address メッセージ
- **国**: 2文字の国コード（必須）
- **都道府県**: 1-50文字（必須）
- **市区町村**: 1-100文字（必須）
- **番地**: 最大200文字（オプション）
- **郵便番号**: 7桁の数字（オプション）

## セットアップ

### 前提条件

- Go 1.21以上
- Protocol Buffersコンパイラ（protoc）
- Make

### インストール

1. リポジトリをクローンまたはダウンロード

2. Protocol Buffersコンパイラをインストール：

```bash
# Ubuntu/Debian (WSL含む)
sudo apt update
sudo apt install protobuf-compiler

# インストール確認
protoc --version
```

**もし上記でエラーが出る場合（パッケージリポジトリの問題）：**

```bash
# パッケージキャッシュをクリア
sudo apt clean
sudo apt autoclean

# 再度試行
sudo apt update
sudo apt install protobuf-compiler
```

**または手動インストール：**

```bash
# 最新版をダウンロード
wget https://github.com/protocolbuffers/protobuf/releases/download/v21.12/protoc-21.12-linux-x86_64.zip
unzip protoc-21.12-linux-x86_64.zip
sudo cp bin/protoc /usr/local/bin/
sudo cp -r include/* /usr/local/include/

# インストール確認
protoc --version
```

3. 依存関係をインストール：

```bash
make install-deps
```

## 使用方法

### コード生成

Protocol BuffersファイルからGoコードを生成：

```bash
make generate
```

### サンプルプログラムの実行

```bash
make run
```

または

```bash
go run main.go
```

### テストの実行

```bash
make test
```

### 全タスクの実行

```bash
make all
```

### クリーンアップ

生成されたファイルを削除：

```bash
make clean
```

## 利用可能なMakeコマンド

| コマンド | 説明 |
|---------|------|
| `make install-deps` | 依存関係をインストール |
| `make generate` | protoファイルからGoコードを生成 |
| `make clean` | 生成されたファイルを削除 |
| `make run` | サンプルプログラムを実行 |
| `make test` | テストを実行 |
| `make all` | 全てのタスクを実行 |
| `make help` | ヘルプを表示 |

## サンプルコードの説明

### メインプログラム（main.go）

メインプログラムでは以下の機能をデモンストレーションします：

1. **有効なユーザーデータのバリデーション**
2. **無効なユーザーデータのバリデーション**
3. **バリデーション関数の使用例**

### テストコード（user_test.go）

包括的なテストケースを含みます：

- 有効なデータのテスト
- 無効なデータのテスト
- エッジケースのテスト

## バリデーションルールの例

### 基本的なバリデーション

```protobuf
message User {
  int64 id = 1 [(validate.rules).int64.gt = 0];
  string name = 2 [(validate.rules).string = {
    min_len: 3,
    max_len: 20,
    pattern: "^[a-zA-Z0-9_]+$"
  }];
  string email = 3 [(validate.rules).string.email = true];
}
```

### 複雑なバリデーション

```protobuf
message User {
  repeated string tags = 7 [(validate.rules).repeated = {
    items: {
      string: {
        min_len: 1,
        max_len: 10
      }
    },
    max_items: 5
  }];
}
```

## 生成されるコードの使用例

```go
user := &user.User{
    Id:    12345,
    Name:  "john_doe",
    Email: "john.doe@example.com",
    Age:   25,
}

// バリデーション実行
if err := user.Validate(); err != nil {
    log.Printf("バリデーションエラー: %v", err)
    return
}

// バリデーション成功
fmt.Println("ユーザーデータは有効です")
```

## トラブルシューティング

### よくある問題

1. **protocが見つからない**
   ```bash
   # Ubuntu/Debian (WSL含む)
   sudo apt update
   sudo apt install protobuf-compiler
   
   # インストール確認
   protoc --version
   ```
   - インストール後もエラーが出る場合は、PATHを確認してください

2. **依存関係のインストールエラー**
   - Goのバージョンが1.21以上であることを確認してください
   - `go mod download`を実行してください

3. **コード生成エラー**
   - `make install-deps`を実行して必要なツールをインストールしてください
   - protoファイルの構文を確認してください

## 参考資料

- [proto-gen-validate公式リポジトリ](https://github.com/envoyproxy/protoc-gen-validate)
- [Protocol Buffers公式ドキュメント](https://developers.google.com/protocol-buffers)
- [Go言語公式ドキュメント](https://golang.org/doc/)

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
