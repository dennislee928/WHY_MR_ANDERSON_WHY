#!/bin/bash

# =============================================================================
# Pandora Box Console IDS-IPS - PaaS 多平台部署腳本
# =============================================================================
# 此腳本自動化部署到多個 PaaS 平台
# 作者: Pandora Security Team
# 版本: 1.0.0
# =============================================================================

set -e  # 遇到錯誤立即退出

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日誌函數
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 顯示標題
show_header() {
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║   Pandora Box Console IDS-IPS - PaaS 多平台部署工具          ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
}

# 檢查必要工具
check_prerequisites() {
    log_info "檢查必要工具..."
    
    local missing_tools=()
    
    # 檢查 Docker
    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi
    
    # 檢查 Git
    if ! command -v git &> /dev/null; then
        missing_tools+=("git")
    fi
    
    # 檢查各平台 CLI
    if ! command -v railway &> /dev/null; then
        log_warning "Railway CLI 未安裝，將跳過 Railway 部署"
    fi
    
    if ! command -v render &> /dev/null; then
        log_warning "Render CLI 未安裝，將跳過 Render 部署"
    fi
    
    if ! command -v koyeb &> /dev/null; then
        log_warning "Koyeb CLI 未安裝，將跳過 Koyeb 部署"
    fi
    
    if ! command -v fly &> /dev/null; then
        log_warning "Fly CLI 未安裝，將跳過 Fly.io 部署"
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        log_error "缺少必要工具: ${missing_tools[*]}"
        log_info "請安裝缺少的工具後再次執行"
        exit 1
    fi
    
    log_success "所有必要工具已安裝"
}

# 載入環境變數
load_env() {
    log_info "載入環境變數..."
    
    if [ -f ".env.paas" ]; then
        export $(grep -v '^#' .env.paas | xargs)
        log_success "環境變數已載入"
    else
        log_warning ".env.paas 檔案不存在"
        log_info "請從 env.paas.example 複製並設定環境變數"
        log_info "執行: cp env.paas.example .env.paas"
        exit 1
    fi
}

# 部署 Railway (PostgreSQL)
deploy_railway() {
    log_info "開始部署 Railway (PostgreSQL)..."
    
    if ! command -v railway &> /dev/null; then
        log_warning "跳過 Railway 部署 (CLI 未安裝)"
        return
    fi
    
    log_info "1. 初始化 Railway 專案..."
    railway init --name pandora-postgresql || true
    
    log_info "2. 添加 PostgreSQL 資料庫..."
    railway add --database postgresql || true
    
    log_info "3. 上傳初始化 SQL..."
    railway up configs/postgres/init.sql || true
    
    log_success "Railway (PostgreSQL) 部署完成"
    log_info "請在 Railway Dashboard 複製 DATABASE_URL 到 .env.paas"
}

# 部署 Render (Redis + Nginx)
deploy_render() {
    log_info "開始部署 Render (Redis + Nginx)..."
    
    if ! command -v render &> /dev/null; then
        log_warning "跳過 Render 部署 (CLI 未安裝)"
        log_info "請手動在 Render Dashboard 部署:"
        log_info "1. 新增 Redis 資料庫"
        log_info "2. 新增 Web Service (Nginx) 使用 render.yaml"
        return
    fi
    
    log_info "1. 部署 Render 服務..."
    render deploy --yes || true
    
    log_success "Render (Redis + Nginx) 部署完成"
}

# 部署 Koyeb (pandora-agent)
deploy_koyeb() {
    log_info "開始部署 Koyeb (Pandora Agent + Promtail)..."
    
    if ! command -v koyeb &> /dev/null; then
        log_warning "跳過 Koyeb 部署 (CLI 未安裝)"
        log_info "請手動在 Koyeb Dashboard 部署:"
        log_info "1. 建立新 App"
        log_info "2. 使用 Dockerfile.agent.koyeb"
        log_info "3. 設定環境變數"
        return
    fi
    
    log_info "1. 建立 Koyeb 應用..."
    koyeb app create pandora-agent || true
    
    log_info "2. 建置並推送 Docker 映像..."
    docker build -f Dockerfile.agent.koyeb -t pandora-agent:latest .
    
    log_info "3. 部署到 Koyeb..."
    koyeb service create pandora-agent \
        --app pandora-agent \
        --docker pandora-agent:latest \
        --ports 8080:http \
        --routes /:8080 \
        --regions fra \
        --instance-type nano \
        --replicas 2 || true
    
    log_success "Koyeb (Pandora Agent) 部署完成"
}

# 部署 Patr.io (axiom-ui)
deploy_patr() {
    log_info "開始部署 Patr.io (Axiom UI)..."
    
    log_warning "Patr.io CLI 功能有限"
    log_info "請手動在 Patr.io Dashboard 部署:"
    log_info "1. 建立新 Container"
    log_info "2. 連接 GitHub Repository"
    log_info "3. 使用 patr.yaml 配置"
    log_info "4. 設定環境變數"
    log_info "5. 部署應用"
}

# 部署 Fly.io (監控系統)
deploy_flyio() {
    log_info "開始部署 Fly.io (監控系統)..."
    
    if ! command -v fly &> /dev/null; then
        log_warning "跳過 Fly.io 部署 (CLI 未安裝)"
        log_info "請安裝 Fly CLI: curl -L https://fly.io/install.sh | sh"
        return
    fi
    
    log_info "1. 建立 Fly.io 應用..."
    fly apps create pandora-monitoring --org personal || true
    
    log_info "2. 建立持久化儲存..."
    fly volumes create prometheus_data --size 3 --region nrt || true
    fly volumes create loki_data --size 3 --region nrt || true
    fly volumes create grafana_data --size 1 --region nrt || true
    fly volumes create alertmanager_data --size 1 --region nrt || true
    
    log_info "3. 設定環境變數..."
    fly secrets set \
        GRAFANA_ADMIN_PASSWORD="${GRAFANA_ADMIN_PASSWORD}" \
        LOG_LEVEL=info || true
    
    log_info "4. 部署應用..."
    fly deploy --config fly.toml --dockerfile Dockerfile.monitoring
    
    log_success "Fly.io (監控系統) 部署完成"
    log_info "Grafana URL: https://pandora-monitoring.fly.dev"
}

# 顯示部署摘要
show_deployment_summary() {
    echo ""
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                      部署摘要                                  ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo ""
    echo "🎯 部署平台配置："
    echo ""
    echo "📊 Railway.app (PostgreSQL)"
    echo "   - 狀態: 請檢查 Railway Dashboard"
    echo "   - URL: https://railway.app/project/xxx"
    echo ""
    echo "🗄️  Render (Redis + Nginx)"
    echo "   - Redis URL: ${RENDER_REDIS_URL:-未設定}"
    echo "   - Nginx URL: ${RENDER_NGINX_URL:-未設定}"
    echo ""
    echo "🚀 Koyeb (Pandora Agent)"
    echo "   - Agent URL: ${KOYEB_AGENT_URL:-未設定}"
    echo ""
    echo "🖥️  Patr.io (Axiom UI)"
    echo "   - UI URL: ${PATR_UI_URL:-未設定}"
    echo ""
    echo "📈 Fly.io (監控系統)"
    echo "   - Grafana: ${GRAFANA_URL:-未設定}"
    echo "   - Prometheus: ${PROMETHEUS_URL:-未設定}"
    echo "   - Loki: ${LOKI_URL:-未設定}"
    echo "   - AlertManager: ${ALERTMANAGER_URL:-未設定}"
    echo ""
    echo "📝 下一步操作："
    echo "1. 確認所有服務已成功部署"
    echo "2. 更新 .env.paas 中的實際 URL"
    echo "3. 執行驗證腳本: ./scripts/verify-paas-deployment.sh"
    echo "4. 檢查各平台的 Dashboard 和日誌"
    echo ""
}

# 主函數
main() {
    show_header
    
    log_info "開始 PaaS 多平台部署流程..."
    echo ""
    
    # 檢查前置條件
    check_prerequisites
    echo ""
    
    # 載入環境變數
    load_env
    echo ""
    
    # 詢問使用者要部署哪些平台
    echo "請選擇要部署的平台（可多選，用空格分隔）："
    echo "1) Railway (PostgreSQL)"
    echo "2) Render (Redis + Nginx)"
    echo "3) Koyeb (Pandora Agent)"
    echo "4) Patr.io (Axiom UI)"
    echo "5) Fly.io (監控系統)"
    echo "6) 全部部署"
    echo ""
    read -p "請輸入選項 (例如: 1 2 3): " -a choices
    
    echo ""
    
    for choice in "${choices[@]}"; do
        case $choice in
            1)
                deploy_railway
                echo ""
                ;;
            2)
                deploy_render
                echo ""
                ;;
            3)
                deploy_koyeb
                echo ""
                ;;
            4)
                deploy_patr
                echo ""
                ;;
            5)
                deploy_flyio
                echo ""
                ;;
            6)
                deploy_railway
                echo ""
                deploy_render
                echo ""
                deploy_koyeb
                echo ""
                deploy_patr
                echo ""
                deploy_flyio
                echo ""
                ;;
            *)
                log_warning "無效的選項: $choice"
                ;;
        esac
    done
    
    # 顯示部署摘要
    show_deployment_summary
    
    log_success "PaaS 多平台部署流程完成！"
}

# 執行主函數
main "$@"

