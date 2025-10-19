# Pandora Cyber AI/Quantum Security Container
# 基於 Python 3.11 的網路安全 AI/ML 和量子密碼學容器

FROM python:3.11-slim AS builder

# 安裝編譯依賴
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    libopenblas-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# 安裝 Python 依賴
COPY Experimental/cyber-ai-quantum/requirements.txt .
RUN pip install --no-cache-dir --user -r requirements.txt

# ========== 最終階段 ==========

FROM python:3.11-slim

# 安裝運行時依賴
RUN apt-get update && apt-get install -y \
    libopenblas0 \
    curl \
    && rm -rf /var/lib/apt/lists/*

# 創建應用用戶
RUN groupadd -r pandora && useradd -r -g pandora pandora

# 創建目錄
RUN mkdir -p /app/models /app/logs /app/data /app/results /app/qasm_output /app/data/windows_logs && \
    chown -R pandora:pandora /app

WORKDIR /app

# 從 builder 複製 Python 套件到標準位置
COPY --from=builder /root/.local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=builder /root/.local/bin /usr/local/bin

# 複製應用程式核心模組
COPY Experimental/cyber-ai-quantum/ml_threat_detector.py .
COPY Experimental/cyber-ai-quantum/quantum_crypto_sim.py .
COPY Experimental/cyber-ai-quantum/ai_governance.py .
COPY Experimental/cyber-ai-quantum/dataflow_monitor.py .
COPY Experimental/cyber-ai-quantum/advanced_quantum_features.py .
COPY Experimental/cyber-ai-quantum/zero_trust_context.py .
COPY Experimental/cyber-ai-quantum/quantum_ml_hybrid.py .
COPY Experimental/cyber-ai-quantum/quantum_rl_optimizer.py .
COPY Experimental/cyber-ai-quantum/advanced_quantum_algorithms.py .
COPY Experimental/cyber-ai-quantum/main.py .

# 複製新增的量子機器學習模組
COPY Experimental/cyber-ai-quantum/feature_extractor.py .
COPY Experimental/cyber-ai-quantum/generate_dynamic_qasm.py .
COPY Experimental/cyber-ai-quantum/analyze_results.py .
COPY Experimental/cyber-ai-quantum/daily_quantum_job.py .
COPY Experimental/cyber-ai-quantum/train_quantum_classifier.py .

# 複製服務模組
COPY Experimental/cyber-ai-quantum/services/ ./services/

# 設定環境變數
ENV PYTHONUNBUFFERED=1 \
    HOST=0.0.0.0 \
    PORT=8000 \
    LOG_LEVEL=info \
    ZERO_TRUST_ENABLED=true \
    QUANTUM_HYBRID_MODE=true \
    AUTO_RESPONSE_ENABLED=false \
    PREDICTION_THRESHOLD=0.7 \
    RL_LEARNING_RATE=0.1 \
    QUANTUM_DIMENSIONS=10 \
    IBM_QUANTUM_TOKEN="" \
    QUANTUM_BACKEND="local_simulator"

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8000/health || exit 1

# 暴露端口
EXPOSE 8000

# 切換到應用用戶
USER pandora

# 啟動服務
CMD ["python", "main.py"]

