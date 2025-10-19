#!/bin/bash
# Pandora Cyber AI/Quantum - 完整服務測試腳本

echo "=== Pandora Cyber AI/Quantum 服務測試 ==="
echo ""

# 顏色定義
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 測試計數器
TOTAL=0
PASSED=0
FAILED=0

# 測試函數
test_endpoint() {
    local name=$1
    local url=$2
    local method=${3:-GET}
    local data=${4:-}
    
    TOTAL=$((TOTAL + 1))
    echo -n "測試 $name... "
    
    if [ "$method" = "POST" ]; then
        response=$(curl -s -w "\n%{http_code}" -X POST "$url" \
            -H "Content-Type: application/json" \
            -d "$data" 2>/dev/null)
    else
        response=$(curl -s -w "\n%{http_code}" "$url" 2>/dev/null)
    fi
    
    status_code=$(echo "$response" | tail -n1)
    
    if [ "$status_code" = "200" ]; then
        echo -e "${GREEN}✅ 通過${NC} (HTTP $status_code)"
        PASSED=$((PASSED + 1))
    else
        echo -e "${RED}❌ 失敗${NC} (HTTP $status_code)"
        FAILED=$((FAILED + 1))
    fi
}

echo "--- 健康檢查 ---"
test_endpoint "Health Check" "http://localhost:8000/health"
test_endpoint "Root Endpoint" "http://localhost:8000/"
echo ""

echo "--- ML 威脅檢測 ---"
test_endpoint "ML Detect" "http://localhost:8000/api/v1/ml/detect" "POST" \
    '{"source_ip":"192.168.1.100","packets_per_second":1000,"syn_count":50}'
test_endpoint "ML Model Status" "http://localhost:8000/api/v1/ml/model/status"
echo ""

echo "--- 量子密碼學 ---"
test_endpoint "Quantum QKD" "http://localhost:8000/api/v1/quantum/qkd/generate" "POST" \
    '{"key_length":256}'
test_endpoint "Quantum Encrypt" "http://localhost:8000/api/v1/quantum/encrypt" "POST" \
    '{"message":"Test Message"}'
test_endpoint "Quantum Predict" "http://localhost:8000/api/v1/quantum/predict" "POST" \
    '{"historical_threats":[{"severity":0.8,"frequency":0.6,"impact":0.7}]}'
echo ""

echo "--- AI 治理 ---"
test_endpoint "Governance Integrity" "http://localhost:8000/api/v1/governance/integrity"
test_endpoint "Adversarial Detect" "http://localhost:8000/api/v1/governance/adversarial/detect" "POST" \
    '{"source_ip":"192.168.1.100","packets_per_second":100}'
test_endpoint "Governance Report" "http://localhost:8000/api/v1/governance/report"
echo ""

echo "--- 資料流監控 ---"
test_endpoint "DataFlow Stats" "http://localhost:8000/api/v1/dataflow/stats"
test_endpoint "DataFlow Anomalies" "http://localhost:8000/api/v1/dataflow/anomalies"
test_endpoint "DataFlow Baseline" "http://localhost:8000/api/v1/dataflow/baseline"
echo ""

echo "--- 系統狀態 ---"
test_endpoint "System Status" "http://localhost:8000/api/v1/status"
echo ""

# 測試其他服務
echo "--- 相關服務檢查 ---"
test_endpoint "Axiom UI" "http://localhost:3001/"
test_endpoint "Axiom API" "http://localhost:3001/api/v1/status"
test_endpoint "RabbitMQ Mgmt" "http://localhost:15672/"
test_endpoint "Grafana" "http://localhost:3000/api/health"
test_endpoint "Prometheus" "http://localhost:9090/-/healthy"
echo ""

# 總結
echo "========================================"
echo "測試總結："
echo -e "  總計: $TOTAL"
echo -e "  ${GREEN}通過: $PASSED${NC}"
echo -e "  ${RED}失敗: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 所有測試通過！${NC}"
    exit 0
else
    success_rate=$(awk "BEGIN {printf \"%.1f\", ($PASSED/$TOTAL)*100}")
    echo -e "${YELLOW}⚠️  成功率: $success_rate%${NC}"
    exit 1
fi

