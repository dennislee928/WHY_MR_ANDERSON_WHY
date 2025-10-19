# Pandora Box Console IDS-IPS - 根目錄 Makefile
# 
# 用途：整體專案管理、Docker 操作、部署管理
# 
# 注意：
# - 此 Makefile 用於整體專案管理
# - 如需編譯 Go 程式，請使用: cd Application/be && make all
# - 如需開發前端，請使用: cd Application/Fe && npm run dev

# 變數定義
GO_VERSION := 1.24
PROJECT_NAME := pandora-box-console-ids-ips
DOCKER_REGISTRY := your-registry.com
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse HEAD)

# Go 編譯參數
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"
GOFLAGS := -v

# Docker 相關
DOCKER_COMPOSE_FILE := deployments/docker-compose/docker-compose.yml
DOCKER_COMPOSE_TEST_FILE := deployments/docker-compose/docker-compose.test.yml

.PHONY: help build clean test docker-build docker-push deploy

# 預設目標
help: ## 顯示幫助訊息
	@echo "╔════════════════════════════════════════════════════════════╗"
	@echo "║  Pandora Box Console IDS-IPS - 根目錄 Makefile            ║"
	@echo "║  用途：整體專案管理、Docker、部署                          ║"
	@echo "╚════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "📝 注意："
	@echo "  • 編譯 Go 程式: cd Application/be && make all"
	@echo "  • 開發前端:    cd Application/Fe && npm run dev"
	@echo "  • 本地構建:    cd Application && ./build-local.*"
	@echo ""
	@echo "🔧 可用的命令："
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 建置相關（建議使用 Application/be/Makefile）
build: ## 建置所有執行檔（建議：cd Application/be && make all）
	@echo "⚠️  建議使用: cd Application/be && make all"
	@echo "使用根目錄 Makefile 構建..."
	@cd Application/be && make all

build-app: ## 使用 Application 構建腳本
	@echo "使用 Application/build-local 腳本..."
	@cd Application && ./build-local.sh || .\build-local.ps1

# 清理相關
clean: ## 清理建置檔案
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf dist/
	@go clean -cache
	@go clean -testcache

# 測試相關
test: ## 執行單元測試
	@echo "Running unit tests..."
	@go test -v ./...

test-coverage: ## 執行測試並產生覆蓋率報告
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-integration: ## 執行整合測試
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...

test-benchmark: ## 執行效能測試
	@echo "Running benchmark tests..."
	@go test -bench=. -benchmem ./...

# Docker 相關
docker-build: ## 建置 Docker 映像
	@echo "Building Docker images..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build

docker-build-no-cache: ## 重新建置 Docker 映像 (無快取)
	@echo "Building Docker images without cache..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build --no-cache

docker-push: ## 推送 Docker 映像到註冊表
	@echo "Pushing Docker images..."
	@docker tag pandora-agent:latest $(DOCKER_REGISTRY)/pandora-agent:$(VERSION)
	@docker tag axiom-ui:latest $(DOCKER_REGISTRY)/axiom-ui:$(VERSION)
	@docker push $(DOCKER_REGISTRY)/pandora-agent:$(VERSION)
	@docker push $(DOCKER_REGISTRY)/axiom-ui:$(VERSION)

# 部署相關
deploy: ## 部署服務
	@echo "Deploying services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

deploy-dev: ## 部署開發環境
	@echo "Deploying development environment..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) -f docker-compose.dev.yml up -d

stop: ## 停止服務
	@echo "Stopping services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

restart: ## 重啟服務
	@echo "Restarting services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) restart

logs: ## 檢視服務日誌
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# 開發相關
dev-setup: ## 設置開發環境
	@echo "Setting up development environment..."
	@go mod download
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@cp env.example .env

lint: ## 執行程式碼檢查
	@echo "Running linter..."
	@golangci-lint run

fmt: ## 格式化程式碼
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

vet: ## 執行 go vet
	@echo "Running go vet..."
	@go vet ./...

# 文件相關
docs: ## 產生 API 文件
	@echo "Generating API documentation..."
	@swag init -g cmd/ui/main.go -o docs/

docs-serve: ## 啟動文件伺服器
	@echo "Starting documentation server..."
	@cd docs && python3 -m http.server 8000

# 監控相關
metrics: ## 檢視指標
	@echo "Fetching metrics..."
	@curl -s http://localhost:8080/metrics

health: ## 檢查服務健康狀態
	@echo "Checking service health..."
	@curl -s http://localhost:8080/health | jq .

status: ## 檢視服務狀態
	@echo "Service status:"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) ps

# 資料庫相關
db-migrate: ## 執行資料庫遷移
	@echo "Running database migrations..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec postgres psql -U pandora -d pandora -f /migrations/init.sql

db-backup: ## 備份資料庫
	@echo "Backing up database..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec postgres pg_dump -U pandora pandora > backup_$(shell date +%Y%m%d_%H%M%S).sql

db-restore: ## 還原資料庫 (需要指定檔案: make db-restore FILE=backup.sql)
	@echo "Restoring database from $(FILE)..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T postgres psql -U pandora -d pandora < $(FILE)

# 安全相關
security-scan: ## 執行安全掃描
	@echo "Running security scan..."
	@go list -json -m all | nancy sleuth

vulnerability-check: ## 檢查已知漏洞
	@echo "Checking for vulnerabilities..."
	@govulncheck ./...

# 效能相關
profile-cpu: ## 產生 CPU 效能分析
	@echo "Generating CPU profile..."
	@go tool pprof -http=:8080 http://localhost:8080/debug/pprof/profile

profile-memory: ## 產生記憶體效能分析
	@echo "Generating memory profile..."
	@go tool pprof -http=:8080 http://localhost:8080/debug/pprof/heap

# 版本相關
version: ## 顯示版本資訊
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"

release: ## 建立發布版本
	@echo "Creating release $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)

# 工具相關
tools: ## 安裝開發工具
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/sonatard/noctx/cmd/noctx@latest

# 快速啟動
quick-start: dev-setup build deploy ## 快速啟動 (設置+建置+部署)

# 完整測試
full-test: lint vet test test-integration ## 執行完整測試套件
