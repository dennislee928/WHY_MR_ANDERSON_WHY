FROM python:3.11-slim

# 安裝必要的系統依賴
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 安裝 Robot Framework 和必要套件
RUN pip install --no-cache-dir \
    robotframework>=6.0.0 \
    robotframework-requests>=0.9.4 \
    robotframework-jsonlibrary>=0.5 \
    requests>=2.28.0

# 複製 entrypoint 腳本
COPY entrypoint.sh /entrypoint.sh

# 設定執行權限
RUN chmod +x /entrypoint.sh

# 設定工作目錄
WORKDIR /workspace

# 設定 entrypoint
ENTRYPOINT ["/entrypoint.sh"]


