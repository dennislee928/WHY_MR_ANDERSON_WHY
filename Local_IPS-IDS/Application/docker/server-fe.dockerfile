# Axiom UI Server Dockerfile
FROM golang:1.24-alpine AS builder

# 安裝依賴
RUN apk add --no-cache git ca-certificates tzdata

# 設定工作目錄
WORKDIR /app

# 複製 go mod 檔案
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製原始碼
COPY . .

# 建置 UI 伺服器
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o axiom-ui ./cmd/ui

# 最終階段
FROM alpine:latest

# 安裝必要工具
RUN apk --no-cache add ca-certificates curl wget

# 建立非 root 使用者
RUN addgroup -g 1001 -S axiom && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G axiom -g axiom axiom

# 設定工作目錄
WORKDIR /app

# 從建置階段複製執行檔
COPY --from=builder /app/axiom-ui .

# 建立必要目錄
RUN mkdir -p /app/web /app/data /app/logs && \
    chown -R axiom:axiom /app

# 複製前端資源
COPY --chown=axiom:axiom web/ ./web/

# 複製設定檔範本
COPY --chown=axiom:axiom configs/ui-config.yaml.template ./ui-config.yaml

# 設定權限
RUN chmod +x axiom-ui

# 暴露端口
EXPOSE 3001

# 使用非 root 使用者執行
USER axiom

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:3001/api/v1/status || exit 1

# 啟動命令
CMD ["./axiom-ui", "--config", "/app/ui-config.yaml"]
