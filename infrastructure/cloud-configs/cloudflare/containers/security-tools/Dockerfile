# Security Tools Container
FROM alpine:latest

WORKDIR /app

# Install security tools
RUN apk add --no-cache \
    nmap \
    nikto \
    sqlmap \
    gobuster \
    dirb \
    curl \
    jq \
    python3 \
    py3-pip

# Install additional Python tools
RUN pip3 install --no-cache-dir \
    requests \
    beautifulsoup4 \
    scapy

# Copy security tools scripts
COPY . .

# Create non-root user
RUN adduser -D -s /bin/sh security
RUN chown -R security:security /app
USER security

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Start security tools API
CMD ["python3", "security_api.py"]
