# proto-gen-validate サンプルプロジェクトのMakefile

# 変数定義
PROTO_DIR := .
GO_OUT_DIR := .
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
PROTOC_GEN_VALIDATE := protoc-gen-validate

# Go modulesのパスを取得
VALIDATE_PROTO_PATH := $(shell go list -m -f '{{.Dir}}' github.com/envoyproxy/protoc-gen-validate)

# 依存関係のインストール
.PHONY: install-deps
install-deps:
	@echo "依存関係をインストール中..."
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

# protoファイルからGoコードを生成
.PHONY: generate
generate: install-deps
	@echo "protoファイルからGoコードを生成中..."
	$(PROTOC) \
		--proto_path=. \
		--proto_path=$(VALIDATE_PROTO_PATH) \
		--go_out=$(GO_OUT_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) \
		--go-grpc_opt=paths=source_relative \
		--validate_out="lang=go:$(GO_OUT_DIR)" \
		user.proto
	@echo "コード生成完了！"

# 生成されたファイルをクリーンアップ
.PHONY: clean
clean:
	@echo "生成されたファイルを削除中..."
	rm -f user/*.pb.go
	rm -f user/*_validate.pb.go
	@echo "クリーンアップ完了！"

# サンプルプログラムを実行
.PHONY: run
run: generate
	@echo "サンプルプログラムを実行中..."
	go run main.go

# テストを実行
.PHONY: test
test: generate
	@echo "テストを実行中..."
	go test ./...

# 全てのタスクを実行（依存関係インストール + コード生成 + 実行）
.PHONY: all
all: generate run

# ヘルプを表示
.PHONY: help
help:
	@echo "利用可能なコマンド:"
	@echo "  install-deps  - 依存関係をインストール"
	@echo "  generate      - protoファイルからGoコードを生成"
	@echo "  clean         - 生成されたファイルを削除"
	@echo "  run           - サンプルプログラムを実行"
	@echo "  test          - テストを実行"
	@echo "  all           - 全てのタスクを実行"
	@echo "  help          - このヘルプを表示"
