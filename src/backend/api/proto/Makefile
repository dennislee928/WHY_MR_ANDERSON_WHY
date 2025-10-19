# Makefile for generating gRPC code from proto files
# 用於從 proto 文件生成 gRPC 代碼

.PHONY: all clean install generate device network control common

# 默認目標：生成所有 proto 文件
all: generate

# 安裝必要的工具
install:
	@echo "Installing protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Plugins installed"

# 生成所有 gRPC 代碼
generate: common device network control
	@echo "✅ All proto files generated"

# 生成 common.proto
common:
	@echo "Generating common.proto..."
	protoc --go_out=../.. \
		--go_opt=paths=source_relative \
		--go-grpc_out=../.. \
		--go-grpc_opt=paths=source_relative \
		common.proto

# 生成 device.proto
device:
	@echo "Generating device.proto..."
	protoc --go_out=../.. \
		--go_opt=paths=source_relative \
		--go-grpc_out=../.. \
		--go-grpc_opt=paths=source_relative \
		device.proto

# 生成 network.proto
network:
	@echo "Generating network.proto..."
	protoc --go_out=../.. \
		--go_opt=paths=source_relative \
		--go-grpc_out=../.. \
		--go-grpc_opt=paths=source_relative \
		network.proto

# 生成 control.proto
control:
	@echo "Generating control.proto..."
	protoc --go_out=../.. \
		--go_opt=paths=source_relative \
		--go-grpc_out=../.. \
		--go-grpc_opt=paths=source_relative \
		control.proto

# 清理生成的文件
clean:
	@echo "Cleaning generated files..."
	rm -f *.pb.go
	rm -f *_grpc.pb.go
	@echo "✅ Cleaned"

# 驗證 proto 文件語法
validate:
	@echo "Validating proto files..."
	protoc --descriptor_set_out=/dev/null common.proto
	protoc --descriptor_set_out=/dev/null device.proto
	protoc --descriptor_set_out=/dev/null network.proto
	protoc --descriptor_set_out=/dev/null control.proto
	@echo "✅ All proto files are valid"

# 顯示幫助
help:
	@echo "Makefile for gRPC code generation"
	@echo ""
	@echo "Usage:"
	@echo "  make install    - Install protoc plugins"
	@echo "  make generate   - Generate all gRPC code"
	@echo "  make device     - Generate device service code"
	@echo "  make network    - Generate network service code"
	@echo "  make control    - Generate control service code"
	@echo "  make common     - Generate common types"
	@echo "  make validate   - Validate proto files"
	@echo "  make clean      - Clean generated files"
	@echo "  make help       - Show this help message"

