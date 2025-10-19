# Koyeb Pandora Agent Dockerfile
# 整合 pandora-agent 與 promtail 於同一容器

FROM golang:1.24-alpine3.21 AS builder

# 安裝建置依賴
RUN apk add --no-cache git gcc musl-dev

# 設定工作目錄
WORKDIR /app

# 複製依賴檔案
COPY go.mod go.sum ./
RUN go mod download

# 複製原始碼
COPY . .

# 建置 Agent
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o pandora-agent ./cmd/agent

# 最終階段
FROM alpine:3.21

# 安裝運行時依賴
RUN apk add --no-cache ca-certificates curl wget unzip supervisor file

# 安裝 Promtail - 調試版本，跳過版本檢查避免阻塞
ARG PROMTAIL_VERSION=2.9.2
RUN wget -q -O /tmp/promtail.zip https://github.com/grafana/loki/releases/download/v${PROMTAIL_VERSION}/promtail-linux-amd64.zip && \
    cd /tmp && \
    unzip promtail.zip && \
    ls -la promtail-linux-amd64 && \
    file promtail-linux-amd64 && \
    chmod +x promtail-linux-amd64 && \
    mv promtail-linux-amd64 /usr/local/bin/promtail && \
    ls -la /usr/local/bin/promtail && \
    rm -f promtail.zip

# 建立非 root 使用者
RUN addgroup -S pandora && adduser -S pandora -G pandora

# 建立必要目錄
RUN mkdir -p /app/data /app/logs /app/configs /etc/supervisor.d && \
    chown -R pandora:pandora /app

# 複製編譯好的執行檔
COPY --from=builder /app/pandora-agent /app/

# 複製配置檔案
COPY configs/agent-config.yaml /app/configs/
COPY configs/promtail-paas.yaml /app/configs/promtail.yaml

# 複製 Supervisor 配置
COPY configs/supervisord-koyeb.conf /etc/supervisord.conf

# 切換到應用目錄
WORKDIR /app

# 暴露端口
EXPOSE 8080

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 切換到非 root 用戶
USER pandora

# 使用 supervisor 同時運行 agent 和 promtail
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]

