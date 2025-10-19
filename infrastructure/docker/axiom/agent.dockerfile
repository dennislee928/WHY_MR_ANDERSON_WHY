# Pandora Box Console IDS-IPS Agent Dockerfile
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

# 建置應用程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pandora-agent ./cmd/agent

# 最終階段
FROM alpine:latest

# 安裝必要工具
RUN apk --no-cache add ca-certificates curl wget netcat-openbsd

# 建立非 root 使用者
RUN addgroup -g 1001 -S pandora && \
    adduser -S -D -H -u 1001 -h /app -s /sbin/nologin -G pandora -g pandora pandora

# 設定工作目錄
WORKDIR /app

# 從建置階段複製執行檔
COPY --from=builder /app/pandora-agent .

# 建立必要目錄
RUN mkdir -p /app/data /app/logs /app/certs && \
    chown -R pandora:pandora /app

# 複製設定檔範本
COPY --chown=pandora:pandora configs/agent-config.yaml.template ./agent-config.yaml

# 設定權限
RUN chmod +x pandora-agent

# 健康檢查腳本
COPY --chown=pandora:pandora scripts/health-check.sh ./
RUN chmod +x health-check.sh

# 暴露端口
EXPOSE 8080 8443

# 使用非 root 使用者執行
USER pandora

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD ./health-check.sh

# 啟動命令
CMD ["./pandora-agent", "--config", "/app/agent-config.yaml"]
