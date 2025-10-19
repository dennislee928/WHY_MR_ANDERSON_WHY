# Patr.io Axiom UI Dockerfile
# 優化的輕量化 UI 伺服器容器

FROM golang:1.24-alpine3.21 AS builder

# 安裝建置依賴
RUN apk add --no-cache git gcc musl-dev nodejs npm

# 設定工作目錄
WORKDIR /app

# 複製依賴檔案
COPY go.mod go.sum ./
RUN go mod download

# 複製原始碼
COPY . .

# 建置前端資源
WORKDIR /app/Application/Fe
RUN if [ -f package.json ]; then npm install && npm run build; fi

# 建置 UI Server
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o axiom-ui ./cmd/ui

# 最終階段
FROM alpine:3.21

# 安裝運行時依賴
RUN apk add --no-cache ca-certificates curl tzdata

# 設定時區
ENV TZ=Asia/Taipei

# 建立非 root 使用者
RUN addgroup -S pandora && adduser -S pandora -G pandora

# 建立必要目錄
RUN mkdir -p /app/web /app/configs /app/data && \
    chown -R pandora:pandora /app

# 複製編譯好的執行檔
COPY --from=builder /app/axiom-ui /app/

# 複製前端資源（簡單 HTML）
COPY --from=builder /app/Application/Fe/public/index.html /app/web/index.html
COPY --from=builder /app/Application/Fe/public/favicon.ico /app/web/favicon.ico

# 複製配置檔案
COPY configs/ui-config.yaml /app/configs/

# 切換到非 root 使用者
USER pandora

# 切換到應用目錄
WORKDIR /app

# 暴露端口
EXPOSE 3001

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3001/api/v1/status || exit 1

# 啟動服務
CMD ["./axiom-ui", "--config", "/app/configs/ui-config.yaml"]

