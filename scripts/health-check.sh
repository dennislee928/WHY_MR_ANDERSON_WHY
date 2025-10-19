#!/bin/bash
# Pandora Box Console Health Check Script

set -e

# 設定變數
HEALTH_ENDPOINT="http://localhost:8080/health"
METRICS_ENDPOINT="http://localhost:8080/metrics"
TIMEOUT=10

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 日誌函數
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 檢查 HTTP 端點
check_endpoint() {
    local url=$1
    local description=$2
    
    if curl -f -s --max-time $TIMEOUT "$url" > /dev/null 2>&1; then
        log_info "$description is healthy"
        return 0
    else
        log_error "$description is not responding"
        return 1
    fi
}

# 檢查裝置連接
check_device() {
    if [ -c "/dev/ttyUSB0" ]; then
        log_info "USB device is connected"
        return 0
    else
        log_warn "USB device not found at /dev/ttyUSB0"
        return 1
    fi
}

# 檢查網路介面
check_network() {
    if ip link show eth0 > /dev/null 2>&1; then
        log_info "Network interface eth0 is available"
        return 0
    else
        log_warn "Network interface eth0 not found"
        return 1
    fi
}

# 檢查記憶體使用率
check_memory() {
    local memory_usage
    memory_usage=$(free | grep Mem | awk '{printf "%.0f", $3/$2 * 100.0}')
    
    if [ "$memory_usage" -lt 90 ]; then
        log_info "Memory usage is ${memory_usage}%"
        return 0
    else
        log_warn "High memory usage: ${memory_usage}%"
        return 1
    fi
}

# 檢查磁碟空間
check_disk() {
    local disk_usage
    disk_usage=$(df / | tail -1 | awk '{print $5}' | sed 's/%//')
    
    if [ "$disk_usage" -lt 90 ]; then
        log_info "Disk usage is ${disk_usage}%"
        return 0
    else
        log_warn "High disk usage: ${disk_usage}%"
        return 1
    fi
}

# 主要健康檢查函數
main() {
    log_info "Starting Pandora Box Console health check..."
    
    local exit_code=0
    
    # 檢查 HTTP 端點
    if ! check_endpoint "$HEALTH_ENDPOINT" "Health endpoint"; then
        exit_code=1
    fi
    
    if ! check_endpoint "$METRICS_ENDPOINT" "Metrics endpoint"; then
        exit_code=1
    fi
    
    # 檢查系統資源
    if ! check_memory; then
        exit_code=1
    fi
    
    if ! check_disk; then
        exit_code=1
    fi
    
    # 檢查硬體
    check_device
    check_network
    
    if [ $exit_code -eq 0 ]; then
        log_info "All health checks passed"
    else
        log_error "Some health checks failed"
    fi
    
    exit $exit_code
}

# 執行主函數
main "$@"
