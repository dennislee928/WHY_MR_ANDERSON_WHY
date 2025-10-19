#!/bin/bash
# PaaS 服務整合驗證腳本
# 用於檢查所有 PaaS 平台上的服務是否正常運行並正確整合

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 服務 URL
KOYEB_AGENT_URL="https://dizzy-sher-mitake-7f13854a.koyeb.app:8080"
FLYIO_MONITORING_URL="https://pandora-monitoring.fly.dev"
RENDER_REDIS_URL="https://redis-7-2-11-alpine3-21.onrender.com"
RENDER_NGINX_URL="https://nginx-stable-perl-boqt.onrender.com"

# 計數器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 測試函數
test_service() {
    local name=$1
    local url=$2
    local expected_code=${3:-200}
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -n "Testing ${name}... "
    
    response=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 10 --max-time 30 "${url}" 2>/dev/null || echo "000")
    
    if [ "$response" -eq "$expected_code" ] || [ "$response" -eq 200 ] || [ "$response" -eq 302 ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP ${response})"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (HTTP ${response})"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
}

# 橫幅
echo -e "${BLUE}"
echo "╔═══════════════════════════════════════════════════════════╗"
echo "║     Pandora Box Console - PaaS Integration Test         ║"
echo "╚═══════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# 1. Koyeb Agent 測試
echo -e "\n${YELLOW}[1/4] Testing Koyeb Agent...${NC}"
test_service "Koyeb Agent Health" "${KOYEB_AGENT_URL}/health"
test_service "Koyeb Agent Metrics" "${KOYEB_AGENT_URL}/metrics"

# 2. Fly.io Monitoring Stack 測試
echo -e "\n${YELLOW}[2/4] Testing Fly.io Monitoring Stack...${NC}"
test_service "Prometheus Health" "${FLYIO_MONITORING_URL}:9090/-/healthy"
test_service "Prometheus Targets" "${FLYIO_MONITORING_URL}:9090/api/v1/targets"
test_service "Grafana Health" "${FLYIO_MONITORING_URL}:3000/api/health"
test_service "Loki Ready" "${FLYIO_MONITORING_URL}:3100/ready"
test_service "AlertManager Health" "${FLYIO_MONITORING_URL}:9093/-/healthy"

# 3. Render Services 測試
echo -e "\n${YELLOW}[3/4] Testing Render Services...${NC}"
test_service "Nginx Proxy" "${RENDER_NGINX_URL}/health" 200

# 4. 整合測試
echo -e "\n${YELLOW}[4/4] Testing Service Integration...${NC}"

# 檢查 Prometheus 是否抓取到 Koyeb Agent
echo -n "Testing Prometheus scraping Koyeb Agent... "
targets=$(curl -s "${FLYIO_MONITORING_URL}:9090/api/v1/targets" 2>/dev/null || echo "")
if echo "$targets" | grep -q "dizzy-sher-mitake"; then
    echo -e "${GREEN}✓ PASS${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${RED}✗ FAIL${NC} (Target not found in Prometheus)"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# 檢查 Grafana 資料源
echo -n "Testing Grafana data sources... "
datasources=$(curl -s -u admin:pandora123 "${FLYIO_MONITORING_URL}:3000/api/datasources" 2>/dev/null || echo "[]")
if echo "$datasources" | grep -q "Prometheus" && echo "$datasources" | grep -q "Loki"; then
    echo -e "${GREEN}✓ PASS${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
else
    echo -e "${YELLOW}⚠ WARN${NC} (Data sources may need configuration)"
    FAILED_TESTS=$((FAILED_TESTS + 1))
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# 結果摘要
echo -e "\n${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}Test Summary${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "Total Tests:  ${TOTAL_TESTS}"
echo -e "${GREEN}Passed:       ${PASSED_TESTS}${NC}"
echo -e "${RED}Failed:       ${FAILED_TESTS}${NC}"
echo -e "Success Rate: $(( PASSED_TESTS * 100 / TOTAL_TESTS ))%"

# 服務狀態表
echo -e "\n${BLUE}Service Status:${NC}"
echo "┌─────────────────────────┬────────────────────────────────────────────────┐"
echo "│ Service                 │ URL                                            │"
echo "├─────────────────────────┼────────────────────────────────────────────────┤"
echo "│ Koyeb Agent             │ ${KOYEB_AGENT_URL}                             │"
echo "│ Fly.io Grafana          │ ${FLYIO_MONITORING_URL}:3000                   │"
echo "│ Fly.io Prometheus       │ ${FLYIO_MONITORING_URL}:9090                   │"
echo "│ Fly.io Loki             │ ${FLYIO_MONITORING_URL}:3100                   │"
echo "│ Fly.io AlertManager     │ ${FLYIO_MONITORING_URL}:9093                   │"
echo "│ Render Nginx            │ ${RENDER_NGINX_URL}                            │"
echo "└─────────────────────────┴────────────────────────────────────────────────┘"

# 建議
if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "\n${YELLOW}⚠ Recommendations:${NC}"
    echo "1. Check service logs for errors"
    echo "2. Verify environment variables are set correctly"
    echo "3. Ensure all services are deployed and running"
    echo "4. Review the integration guide: docs/deployment/paas-integration-guide.md"
    exit 1
else
    echo -e "\n${GREEN}✓ All services are running and integrated successfully!${NC}"
    exit 0
fi
