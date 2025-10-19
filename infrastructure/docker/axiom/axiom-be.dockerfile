# Pandora Box Console - Axiom Backend V3 Dockerfile
# 獨立的 Axiom Backend API 服務 (Go + Gin)
# 版本: 3.1.0 - 支援 Agent 管理、四層儲存、合規性引擎

# ========== 構建階段 ==========
FROM golang:1.23-alpine AS builder

# 安裝構建依賴
RUN apk add --no-cache git make gcc musl-dev

WORKDIR /app

# 複製 Go modules
COPY go.mod go.sum ./
RUN go mod download

# 複製源碼
COPY . .

# 構建 Axiom Backend V3
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-w -s -X main.Version=3.1.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o axiom-be \
    ./cmd/server

# ========== 運行階段 ==========
FROM alpine:latest

# 安裝運行時依賴
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl

# 創建應用用戶
RUN addgroup -g 1000 pandora && \
    adduser -D -u 1000 -G pandora pandora

# 創建目錄
RUN mkdir -p /app/configs /app/logs /app/data /app/certs && \
    chown -R pandora:pandora /app

WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/axiom-be .

# 設置時區
ENV TZ=Asia/Taipei

# 默認環境變量
ENV PORT=3001
ENV LOG_LEVEL=info

# 健康檢查 (V3 使用 /health 端點)
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:3001/health || exit 1

# 暴露端口
EXPOSE 3001

# 切換到應用用戶
USER pandora

# 啟動服務
CMD ["./axiom-be"]

