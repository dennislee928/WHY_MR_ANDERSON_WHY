# 測試執行環境 Dockerfile
FROM golang:1.24-alpine3.21 AS test-runner

# 安裝必要工具
RUN apk add --no-cache \
    curl \
    wget \
    netcat-openbsd \
    git \
    ca-certificates

# 設定工作目錄
WORKDIR /workspace

# 複製 go mod 檔案
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製測試檔案
COPY . .

# 安裝測試工具
RUN go install github.com/stretchr/testify@latest

# 建立測試結果目錄
RUN mkdir -p /test-results

# 設定測試環境變數
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# 創建非 root 用戶
RUN addgroup -g 1000 tester && \
    adduser -D -u 1000 -G tester tester && \
    chown -R tester:tester /workspace /test-results

# 切換到非 root 用戶
USER tester

# 預設執行命令
CMD ["go", "test", "-v", "-coverprofile=/test-results/coverage.out", "./..."]
