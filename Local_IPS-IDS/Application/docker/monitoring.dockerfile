# Fly.io 監控系統整合 Dockerfile
# 包含 Prometheus, Loki, Grafana, AlertManager

FROM alpine:3.21

# 安裝基礎工具
RUN apk add --no-cache \
    ca-certificates \
    curl \
    wget \
    unzip \
    supervisor \
    nginx \
    tzdata

# 設定時區
ENV TZ=Asia/Taipei

# 安裝 Prometheus
ARG PROMETHEUS_VERSION=2.47.0
RUN wget https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz && \
    tar xzf prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz && \
    mv prometheus-${PROMETHEUS_VERSION}.linux-amd64/prometheus /usr/local/bin/ && \
    mv prometheus-${PROMETHEUS_VERSION}.linux-amd64/promtool /usr/local/bin/ && \
    mv prometheus-${PROMETHEUS_VERSION}.linux-amd64/consoles /etc/prometheus/ && \
    mv prometheus-${PROMETHEUS_VERSION}.linux-amd64/console_libraries /etc/prometheus/ && \
    rm -rf prometheus-${PROMETHEUS_VERSION}.linux-amd64*

# 安裝 Loki
ARG LOKI_VERSION=2.9.2
RUN wget https://github.com/grafana/loki/releases/download/v${LOKI_VERSION}/loki-linux-amd64.zip && \
    unzip loki-linux-amd64.zip && \
    mv loki-linux-amd64 /usr/local/bin/loki && \
    chmod +x /usr/local/bin/loki && \
    rm loki-linux-amd64.zip

# 安裝 Grafana
ARG GRAFANA_VERSION=10.2.0
RUN wget https://dl.grafana.com/oss/release/grafana-${GRAFANA_VERSION}.linux-amd64.tar.gz && \
    tar xzf grafana-${GRAFANA_VERSION}.linux-amd64.tar.gz && \
    mv grafana-${GRAFANA_VERSION} /usr/share/grafana && \
    ln -s /usr/share/grafana/bin/grafana-server /usr/local/bin/grafana-server && \
    rm grafana-${GRAFANA_VERSION}.linux-amd64.tar.gz

# 安裝 AlertManager
ARG ALERTMANAGER_VERSION=0.26.0
RUN wget https://github.com/prometheus/alertmanager/releases/download/v${ALERTMANAGER_VERSION}/alertmanager-${ALERTMANAGER_VERSION}.linux-amd64.tar.gz && \
    tar xzf alertmanager-${ALERTMANAGER_VERSION}.linux-amd64.tar.gz && \
    mv alertmanager-${ALERTMANAGER_VERSION}.linux-amd64/alertmanager /usr/local/bin/ && \
    mv alertmanager-${ALERTMANAGER_VERSION}.linux-amd64/amtool /usr/local/bin/ && \
    rm -rf alertmanager-${ALERTMANAGER_VERSION}.linux-amd64*

# 建立必要目錄
RUN mkdir -p \
    /data/prometheus \
    /data/loki \
    /data/grafana \
    /data/alertmanager \
    /etc/prometheus \
    /etc/loki \
    /etc/grafana \
    /etc/alertmanager \
    /etc/nginx/conf.d \
    /var/log/supervisor \
    /var/log/grafana

# 建立符號連結指向統一數據目錄（僅用於 Prometheus, Loki, AlertManager）
RUN ln -s /data/prometheus /prometheus && \
    ln -s /data/loki /loki && \
    ln -s /data/alertmanager /alertmanager

# Grafana 不使用符號連結，直接使用 /data/grafana 目錄
# 透過環境變數指定路徑

# 複製配置檔案
COPY configs/prometheus.yml /etc/prometheus/prometheus.yml
COPY configs/prometheus/rules /etc/prometheus/rules
COPY configs/loki.yaml /etc/loki/loki.yaml
COPY configs/alertmanager.yml /etc/alertmanager/alertmanager.yml
COPY configs/grafana/grafana.ini /etc/grafana/grafana.ini
COPY configs/grafana/provisioning /etc/grafana/provisioning
COPY configs/grafana/dashboards /data/grafana/dashboards

COPY configs/supervisord-flyio.conf /etc/supervisord.conf
COPY configs/nginx/nginx-flyio.conf /etc/nginx/nginx.conf
COPY configs/nginx/monitoring-flyio.conf /etc/nginx/conf.d/default.conf

# 設定權限
RUN chmod -R 777 /data

# 暴露端口
EXPOSE 80 3000 9090 3100 9093

# 健康檢查
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
  CMD curl -f http://localhost:9090/-/healthy && \
      curl -f http://localhost:3100/ready && \
      curl -f http://localhost:3000/api/health && \
      curl -f http://localhost:9093/-/healthy || exit 1

# 切換到非 root 用戶
USER monitoring

# 使用 supervisor 管理所有服務
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]

