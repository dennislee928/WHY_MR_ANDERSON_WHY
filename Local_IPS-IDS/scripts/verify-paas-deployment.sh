#!/bin/bash

# =============================================================================
# Pandora Box Console IDS-IPS - PaaS 部署驗證腳本
# =============================================================================
# 此腳本驗證所有 PaaS 平台的部署狀態
# =============================================================================

set -e

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 日誌函數
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# 顯示標題
show_header() {
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║     Pandora Box Console IDS-IPS - PaaS 部署驗證工具          ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
}

# 載入環境變數
load_env() {
    if [ -f ".env.paas" ]; then
        export $(grep -v '^#' .env.paas | xargs)
    else
        log_error ".env.paas 檔案不存在"
        exit 1
    fi
}

# 檢查 URL 健康狀態
check_health() {
    local name=$1
    local url=$2
    local expected_code=${3:-200}
    
    if [ -z "$url" ]; then
        log_warning "$name: URL 未設定"
        return 1
    fi
    
    log_info "檢查 $name: $url"
    
    local status_code=$(curl -s -o /dev/null -w "%{http_code}" "$url" --max-time 10 || echo "000")
    
    if [ "$status_code" = "$expected_code" ] || [ "$status_code" = "200" ] || [ "$status_code" = "301" ] || [ "$status_code" = "302" ]; then
        log_success "$name 正常運行 (HTTP $status_code)"
        return 0
    else
        log_error "$name 無法訪問 (HTTP $status_code)"
        return 1
    fi
}

# 驗證 Railway PostgreSQL
verify_railway() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "🗄️  Railway.app (PostgreSQL)"
    echo "═══════════════════════════════════════════════════════════════"
    
    if [ -z "$RAILWAY_DATABASE_URL" ]; then
        log_warning "DATABASE_URL 未設定"
        return
    fi
    
    log_info "檢查資料庫連接..."
    
    # 嘗試使用 psql 連接（如果已安裝）
    if command -v psql &> /dev/null; then
        if psql "$RAILWAY_DATABASE_URL" -c "SELECT 1;" &> /dev/null; then
            log_success "PostgreSQL 連接成功"
        else
            log_error "PostgreSQL 連接失敗"
        fi
    else
        log_warning "psql 未安裝，跳過資料庫連接測試"
        log_info "請手動驗證 Railway Dashboard"
    fi
}

# 驗證 Render (Redis + Nginx)
verify_render() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "🗄️  Render (Redis + Nginx)"
    echo "═══════════════════════════════════════════════════════════════"
    
    # 檢查 Redis
    if [ -n "$RENDER_REDIS_URL" ]; then
        log_info "檢查 Redis 連接..."
        if command -v redis-cli &> /dev/null; then
            if redis-cli -u "$RENDER_REDIS_URL" ping &> /dev/null; then
                log_success "Redis 連接成功"
            else
                log_error "Redis 連接失敗"
            fi
        else
            log_warning "redis-cli 未安裝，跳過 Redis 連接測試"
        fi
    else
        log_warning "RENDER_REDIS_URL 未設定"
    fi
    
    # 檢查 Nginx
    check_health "Nginx" "$RENDER_NGINX_URL/health"
}

# 驗證 Koyeb (Pandora Agent)
verify_koyeb() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "🚀 Koyeb (Pandora Agent + Promtail)"
    echo "═══════════════════════════════════════════════════════════════"
    
    check_health "Pandora Agent" "$KOYEB_AGENT_URL/health"
    
    # 檢查 API 端點
    if [ -n "$KOYEB_AGENT_URL" ]; then
        check_health "Agent API" "$KOYEB_AGENT_URL/api/v1/status"
    fi
}

# 驗證 Patr.io (Axiom UI)
verify_patr() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "🖥️  Patr.io (Axiom UI)"
    echo "═══════════════════════════════════════════════════════════════"
    
    check_health "Axiom UI" "$PATR_UI_URL/api/v1/status"
    
    # 檢查 WebSocket 端點
    if [ -n "$PATR_UI_URL" ]; then
        log_info "檢查 WebSocket 端點..."
        log_warning "WebSocket 測試需要專用工具，請手動驗證"
    fi
}

# 驗證 Fly.io (監控系統)
verify_flyio() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "📈 Fly.io (監控系統)"
    echo "═══════════════════════════════════════════════════════════════"
    
    check_health "Prometheus" "$PROMETHEUS_URL/-/healthy"
    check_health "Loki" "$LOKI_URL/ready"
    check_health "Grafana" "$GRAFANA_URL/api/health"
    check_health "AlertManager" "$ALERTMANAGER_URL/-/healthy"
    
    # 檢查 Prometheus 指標
    if [ -n "$PROMETHEUS_URL" ]; then
        log_info "檢查 Prometheus 目標狀態..."
        local targets=$(curl -s "$PROMETHEUS_URL/api/v1/targets" | grep -o '"health":"up"' | wc -l)
        if [ "$targets" -gt 0 ]; then
            log_success "Prometheus 有 $targets 個正常的目標"
        else
            log_warning "Prometheus 沒有正常的目標"
        fi
    fi
}

# 驗證服務間連接
verify_connectivity() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "🔗 服務間連接測試"
    echo "═══════════════════════════════════════════════════════════════"
    
    log_info "測試 Agent 到資料庫的連接..."
    if [ -n "$KOYEB_AGENT_URL" ]; then
        local db_status=$(curl -s "$KOYEB_AGENT_URL/api/v1/health/database" | grep -o '"status":"ok"')
        if [ -n "$db_status" ]; then
            log_success "Agent 到資料庫連接正常"
        else
            log_warning "Agent 到資料庫連接異常"
        fi
    fi
    
    log_info "測試 Agent 到 Redis 的連接..."
    if [ -n "$KOYEB_AGENT_URL" ]; then
        local redis_status=$(curl -s "$KOYEB_AGENT_URL/api/v1/health/redis" | grep -o '"status":"ok"')
        if [ -n "$redis_status" ]; then
            log_success "Agent 到 Redis 連接正常"
        else
            log_warning "Agent 到 Redis 連接異常"
        fi
    fi
}

# 生成驗證報告
generate_report() {
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                      驗證摘要報告                              ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    
    local report_file="paas-deployment-report-$(date +%Y%m%d-%H%M%S).txt"
    
    {
        echo "Pandora Box Console IDS-IPS - PaaS 部署驗證報告"
        echo "生成時間: $(date)"
        echo ""
        echo "部署平台狀態："
        echo ""
        echo "Railway (PostgreSQL): $([ -n "$RAILWAY_DATABASE_URL" ] && echo "已配置" || echo "未配置")"
        echo "Render (Redis): $([ -n "$RENDER_REDIS_URL" ] && echo "已配置" || echo "未配置")"
        echo "Render (Nginx): $([ -n "$RENDER_NGINX_URL" ] && echo "已配置" || echo "未配置")"
        echo "Koyeb (Agent): $([ -n "$KOYEB_AGENT_URL" ] && echo "已配置" || echo "未配置")"
        echo "Patr.io (UI): $([ -n "$PATR_UI_URL" ] && echo "已配置" || echo "未配置")"
        echo "Fly.io (Prometheus): $([ -n "$PROMETHEUS_URL" ] && echo "已配置" || echo "未配置")"
        echo "Fly.io (Loki): $([ -n "$LOKI_URL" ] && echo "已配置" || echo "未配置")"
        echo "Fly.io (Grafana): $([ -n "$GRAFANA_URL" ] && echo "已配置" || echo "未配置")"
        echo "Fly.io (AlertManager): $([ -n "$ALERTMANAGER_URL" ] && echo "已配置" || echo "未配置")"
    } > "$report_file"
    
    log_success "驗證報告已生成: $report_file"
}

# 顯示建議
show_recommendations() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "📋 建議後續操作"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    echo "1. 登入 Grafana 並驗證儀表板："
    echo "   URL: $GRAFANA_URL"
    echo "   用戶: admin"
    echo "   密碼: $GRAFANA_ADMIN_PASSWORD"
    echo ""
    echo "2. 檢查 Prometheus 目標："
    echo "   URL: $PROMETHEUS_URL/targets"
    echo ""
    echo "3. 測試告警規則："
    echo "   URL: $PROMETHEUS_URL/alerts"
    echo ""
    echo "4. 查看日誌聚合："
    echo "   在 Grafana 中使用 Loki 數據源"
    echo ""
    echo "5. 測試 API 端點："
    echo "   curl $KOYEB_AGENT_URL/api/v1/status"
    echo ""
    echo "6. 監控資源使用："
    echo "   - Railway Dashboard: https://railway.app"
    echo "   - Render Dashboard: https://dashboard.render.com"
    echo "   - Koyeb Dashboard: https://app.koyeb.com"
    echo "   - Patr.io Dashboard: https://patr.cloud"
    echo "   - Fly.io Dashboard: https://fly.io/dashboard"
    echo ""
}

# 主函數
main() {
    show_header
    
    log_info "開始驗證 PaaS 部署狀態..."
    echo ""
    
    load_env
    
    verify_railway
    verify_render
    verify_koyeb
    verify_patr
    verify_flyio
    verify_connectivity
    
    generate_report
    show_recommendations
    
    log_success "驗證流程完成！"
}

# 執行主函數
main "$@"

