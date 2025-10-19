# Monitoring Container
FROM prom/prometheus:latest

WORKDIR /app

# Copy configuration files
COPY prometheus.yml /etc/prometheus/prometheus.yml
COPY rules/ /etc/prometheus/rules/

# Copy monitoring scripts
COPY monitoring-scripts/ ./

# Expose ports
EXPOSE 9090 3001

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:9090/-/healthy || exit 1

# Start Prometheus
CMD ["prometheus", "--config.file=/etc/prometheus/prometheus.yml", "--storage.tsdb.path=/prometheus", "--web.console.libraries=/etc/prometheus/console_libraries", "--web.console.templates=/etc/prometheus/consoles", "--web.enable-lifecycle"]
