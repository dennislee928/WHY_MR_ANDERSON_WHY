#!/bin/bash
# Pandora Box Console 測試執行腳本

set -e

# 設定變數
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEST_RESULTS_DIR="$PROJECT_ROOT/test-results"
COVERAGE_FILE="$TEST_RESULTS_DIR/coverage.out"
COVERAGE_HTML="$TEST_RESULTS_DIR/coverage.html"

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 建立測試結果目錄
create_test_dirs() {
    log_step "建立測試結果目錄..."
    mkdir -p "$TEST_RESULTS_DIR"
    log_info "測試結果目錄: $TEST_RESULTS_DIR"
}

# 清理測試環境
cleanup() {
    log_step "清理測試環境..."
    if [ -f docker-compose.test.yml ]; then
        docker-compose -f docker-compose.test.yml down -v --remove-orphans > /dev/null 2>&1 || true
    fi
    log_info "測試環境已清理"
}

# 檢查依賴
check_dependencies() {
    log_step "檢查依賴..."
    
    # 檢查 Go
    if ! command -v go &> /dev/null; then
        log_error "Go 未安裝"
        exit 1
    fi
    log_info "Go 版本: $(go version)"
    
    # 檢查 Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安裝"
        exit 1
    fi
    log_info "Docker 版本: $(docker --version)"
    
    # 檢查 Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose 未安裝"
        exit 1
    fi
    log_info "Docker Compose 版本: $(docker-compose --version)"
}

# 執行單元測試
run_unit_tests() {
    log_step "執行單元測試..."
    
    cd "$PROJECT_ROOT"
    
    # 執行所有單元測試
    go test -v -race -coverprofile="$COVERAGE_FILE" ./internal/... | tee "$TEST_RESULTS_DIR/unit-tests.log"
    
    if [ $? -eq 0 ]; then
        log_info "單元測試通過"
    else
        log_error "單元測試失敗"
        return 1
    fi
    
    # 產生覆蓋率報告
    if [ -f "$COVERAGE_FILE" ]; then
        go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
        log_info "覆蓋率報告已產生: $COVERAGE_HTML"
        
        # 顯示覆蓋率摘要
        COVERAGE_PERCENT=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}')
        log_info "測試覆蓋率: $COVERAGE_PERCENT"
    fi
}

# 執行整合測試
run_integration_tests() {
    log_step "執行整合測試..."
    
    cd "$PROJECT_ROOT"
    
    # 啟動測試環境
    log_info "啟動測試環境..."
    docker-compose -f docker-compose.test.yml up -d
    
    # 等待服務啟動
    log_info "等待服務啟動..."
    sleep 30
    
    # 檢查服務健康狀態
    check_test_services
    
    # 執行整合測試
    log_info "執行整合測試..."
    go test -v -tags=integration ./test/integration/... | tee "$TEST_RESULTS_DIR/integration-tests.log"
    
    local test_result=$?
    
    # 清理測試環境
    log_info "清理測試環境..."
    docker-compose -f docker-compose.test.yml down -v
    
    if [ $test_result -eq 0 ]; then
        log_info "整合測試通過"
        return 0
    else
        log_error "整合測試失敗"
        return 1
    fi
}

# 檢查測試服務狀態
check_test_services() {
    log_step "檢查測試服務狀態..."
    
    local services=(
        "http://localhost:3001/api/v1/status:Axiom UI"
        "http://localhost:8080/health:Pandora Agent"
        "http://localhost:9090/-/healthy:Prometheus"
        "http://localhost:3000/api/health:Grafana"
        "http://localhost:3100/ready:Loki"
    )
    
    for service in "${services[@]}"; do
        IFS=':' read -r url name <<< "$service"
        
        log_info "檢查服務: $name"
        
        local retries=0
        local max_retries=30
        
        while [ $retries -lt $max_retries ]; do
            if curl -f -s "$url" > /dev/null 2>&1; then
                log_info "服務 $name 已就緒"
                break
            fi
            
            retries=$((retries + 1))
            if [ $retries -eq $max_retries ]; then
                log_error "服務 $name 未能在預期時間內啟動"
                docker-compose -f docker-compose.test.yml logs "$name" || true
                return 1
            fi
            
            sleep 2
        done
    done
    
    log_info "所有測試服務已就緒"
}

# 執行效能測試
run_performance_tests() {
    log_step "執行效能測試..."
    
    cd "$PROJECT_ROOT"
    
    # 執行基準測試
    go test -bench=. -benchmem ./internal/... | tee "$TEST_RESULTS_DIR/benchmark.log"
    
    if [ $? -eq 0 ]; then
        log_info "效能測試完成"
    else
        log_error "效能測試失敗"
        return 1
    fi
}

# 執行安全掃描
run_security_scan() {
    log_step "執行安全掃描..."
    
    cd "$PROJECT_ROOT"
    
    # 檢查是否安裝了 gosec
    if command -v gosec &> /dev/null; then
        log_info "執行 gosec 安全掃描..."
        gosec ./... | tee "$TEST_RESULTS_DIR/security-scan.log"
    else
        log_warn "gosec 未安裝，跳過安全掃描"
    fi
    
    # 檢查是否安裝了 govulncheck
    if command -v govulncheck &> /dev/null; then
        log_info "執行漏洞檢查..."
        govulncheck ./... | tee "$TEST_RESULTS_DIR/vulnerability-check.log"
    else
        log_warn "govulncheck 未安裝，跳過漏洞檢查"
    fi
}

# 產生測試報告
generate_test_report() {
    log_step "產生測試報告..."
    
    local report_file="$TEST_RESULTS_DIR/test-report.html"
    
    cat > "$report_file" << EOF
<!DOCTYPE html>
<html>
<head>
    <title>Pandora Box Console 測試報告</title>
    <meta charset="utf-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .header { background: #f5f5f5; padding: 20px; border-radius: 5px; }
        .section { margin: 20px 0; }
        .pass { color: green; }
        .fail { color: red; }
        .warn { color: orange; }
        pre { background: #f8f8f8; padding: 10px; border-radius: 3px; overflow-x: auto; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Pandora Box Console IDS-IPS 測試報告</h1>
        <p>產生時間: $(date)</p>
    </div>
    
    <div class="section">
        <h2>測試摘要</h2>
        <ul>
EOF

    # 檢查各種測試結果
    if [ -f "$TEST_RESULTS_DIR/unit-tests.log" ]; then
        if grep -q "FAIL" "$TEST_RESULTS_DIR/unit-tests.log"; then
            echo "            <li class=\"fail\">單元測試: 失敗</li>" >> "$report_file"
        else
            echo "            <li class=\"pass\">單元測試: 通過</li>" >> "$report_file"
        fi
    fi
    
    if [ -f "$TEST_RESULTS_DIR/integration-tests.log" ]; then
        if grep -q "FAIL" "$TEST_RESULTS_DIR/integration-tests.log"; then
            echo "            <li class=\"fail\">整合測試: 失敗</li>" >> "$report_file"
        else
            echo "            <li class=\"pass\">整合測試: 通過</li>" >> "$report_file"
        fi
    fi
    
    if [ -f "$COVERAGE_FILE" ]; then
        local coverage=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' || echo "未知")
        echo "            <li>測試覆蓋率: $coverage</li>" >> "$report_file"
    fi

    cat >> "$report_file" << EOF
        </ul>
    </div>
    
    <div class="section">
        <h2>測試檔案</h2>
        <ul>
            <li><a href="coverage.html">覆蓋率報告</a></li>
            <li><a href="unit-tests.log">單元測試日誌</a></li>
            <li><a href="integration-tests.log">整合測試日誌</a></li>
            <li><a href="benchmark.log">效能測試日誌</a></li>
        </ul>
    </div>
</body>
</html>
EOF

    log_info "測試報告已產生: $report_file"
}

# 主要函數
main() {
    echo "======================================="
    echo "Pandora Box Console IDS-IPS 測試套件"
    echo "======================================="
    echo ""
    
    # 設定清理陷阱
    trap cleanup EXIT
    
    # 解析參數
    local run_unit=false
    local run_integration=false
    local run_performance=false
    local run_security=false
    local run_all=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --unit)
                run_unit=true
                shift
                ;;
            --integration)
                run_integration=true
                shift
                ;;
            --performance)
                run_performance=true
                shift
                ;;
            --security)
                run_security=true
                shift
                ;;
            --all)
                run_all=true
                shift
                ;;
            -h|--help)
                echo "使用方法: $0 [選項]"
                echo ""
                echo "選項:"
                echo "  --unit          執行單元測試"
                echo "  --integration   執行整合測試"
                echo "  --performance   執行效能測試"
                echo "  --security      執行安全掃描"
                echo "  --all           執行所有測試"
                echo "  -h, --help      顯示幫助"
                exit 0
                ;;
            *)
                log_error "未知參數: $1"
                exit 1
                ;;
        esac
    done
    
    # 如果沒有指定參數，預設執行所有測試
    if [ "$run_unit" = false ] && [ "$run_integration" = false ] && [ "$run_performance" = false ] && [ "$run_security" = false ]; then
        run_all=true
    fi
    
    # 檢查依賴和建立目錄
    check_dependencies
    create_test_dirs
    
    local overall_result=0
    
    # 執行測試
    if [ "$run_all" = true ] || [ "$run_unit" = true ]; then
        if ! run_unit_tests; then
            overall_result=1
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_integration" = true ]; then
        if ! run_integration_tests; then
            overall_result=1
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_performance" = true ]; then
        if ! run_performance_tests; then
            overall_result=1
        fi
    fi
    
    if [ "$run_all" = true ] || [ "$run_security" = true ]; then
        run_security_scan  # 安全掃描不影響整體結果
    fi
    
    # 產生報告
    generate_test_report
    
    # 顯示結果
    echo ""
    if [ $overall_result -eq 0 ]; then
        log_info "所有測試完成！"
        log_info "測試結果目錄: $TEST_RESULTS_DIR"
    else
        log_error "部分測試失敗！"
        log_info "詳細資訊請查看: $TEST_RESULTS_DIR"
    fi
    
    exit $overall_result
}

# 執行主函數
main "$@"
