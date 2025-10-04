# proto-gen-validate サンプルプロジェクトのMakefile

# 変数定義
PROTO_DIR := .
GO_OUT_DIR := user
PROTOC := protoc
PROTOC_GEN_GO := protoc-gen-go
PROTOC_GEN_GO_GRPC := protoc-gen-go-grpc
PROTOC_GEN_VALIDATE := protoc-gen-validate

# validate.protoファイルのパス（GOPATH固定）
VALIDATE_PROTO_PATH := $(GOPATH)/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.2.1
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
generate:
	@echo "protoファイルからGoコードを生成中..."
	@mkdir -p $(GO_OUT_DIR)
	$(PROTOC) \
		-I. \
		-I$(VALIDATE_PROTO_PATH) \
		--go_out=. \
		--go_opt=module=github.com/jun06t/go-sample/pgv \
		--validate_out="lang=go,module=github.com/jun06t/go-sample/pgv:." \
		user.proto
	@echo "コード生成完了！"

# 生成されたファイルをクリーンアップ
.PHONY: clean
clean:
	@echo "生成されたファイルを削除中..."
	rm -f $(GO_OUT_DIR)/*.pb.go
	rm -f $(GO_OUT_DIR)/*_validate.pb.go
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
