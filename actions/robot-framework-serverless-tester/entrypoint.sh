#!/bin/bash
set -euo pipefail

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

# 解析輸入參數
TEST_DIR="${1:-}"
TARGET_URL="${2:-}"
REPORT_DIR="${3:-reports}"
INCLUDE_TAGS="${4:-}"
EXCLUDE_TAGS="${5:-}"
VARIABLE_FILE="${6:-}"
PROCESSES="${7:-1}"
TIMEOUT="${8:-}"
LOG_LEVEL="${9:-INFO}"

# 驗證必要參數
if [ -z "$TEST_DIR" ] || [ -z "$TARGET_URL" ]; then
    log_error "Missing required parameters: test_dir and target_url are required"
    exit 1
fi

# 驗證測試目錄存在
if [ ! -d "$TEST_DIR" ]; then
    log_error "Test directory does not exist: $TEST_DIR"
    exit 1
fi

# 安全性檢查：防止目錄遍歷
if [[ "$TEST_DIR" == *".."* ]] || [[ "$REPORT_DIR" == *".."* ]]; then
    log_error "Invalid path detected. Directory traversal is not allowed."
    exit 1
fi

log_info "Starting Robot Framework E2E Tests"
log_info "Test Directory: $TEST_DIR"
log_info "Target URL: $TARGET_URL"
log_info "Report Directory: $REPORT_DIR"

# 建立報告目錄
mkdir -p "$REPORT_DIR"

# 注入環境變數為 Robot Framework 變數
log_info "Injecting environment variables as Robot Framework variables..."
SECRET_COUNT=0

# 建立變數陣列（避免使用 eval，提高安全性）
ROBOT_VARIABLES=()

# 注入 TARGET_URL 作為 BASE_URL 和 TARGET_URL 變數（常用於測試）
ROBOT_VARIABLES+=("BASE_URL:$TARGET_URL")
ROBOT_VARIABLES+=("TARGET_URL:$TARGET_URL")
SECRET_COUNT=2

# 掃描環境變數，找出可能的 secrets 和配置
# 注意：GitHub Actions 會將 secrets 設為環境變數，但我們只處理明確的變數
# 避免處理系統變數（PATH, HOME, etc.）
EXCLUDE_VARS="PATH|HOME|PWD|SHELL|USER|LANG|LC_|TERM|_|GITHUB_|RUNNER_|INPUT_"

while IFS='=' read -r key value; do
    # 跳過空行
    [[ -z "$key" ]] && continue
    
    # 跳過排除的變數
    if [[ "$key" =~ ^($EXCLUDE_VARS) ]]; then
        continue
    fi
    
    # 只處理大寫變數（通常是配置或 secrets）
    if [[ "$key" =~ ^[A-Z][A-Z0-9_]*$ ]] && [[ -n "${!key:-}" ]]; then
        # Robot Framework 變數值不需要轉義，直接使用
        ROBOT_VARIABLES+=("$key:${!key}")
        SECRET_COUNT=$((SECRET_COUNT + 1))
    fi
done < <(env | sort || true)

log_info "Injected $SECRET_COUNT variables as Robot Framework variables"

# 建構 robot 命令參數陣列（安全執行，避免命令注入）
EXIT_CODE=0
ROBOT_ARGS=()

# 設定基本參數
ROBOT_ARGS+=("--outputdir" "$REPORT_DIR")
ROBOT_ARGS+=("--logdir" "$REPORT_DIR")
ROBOT_ARGS+=("--reportdir" "$REPORT_DIR")
ROBOT_ARGS+=("--loglevel" "$LOG_LEVEL")

# 添加標籤過濾
if [ -n "$INCLUDE_TAGS" ]; then
    IFS=' ' read -ra TAGS <<< "$(echo "$INCLUDE_TAGS" | tr ',' ' ')"
    for tag in "${TAGS[@]}"; do
        ROBOT_ARGS+=("--include" "$tag")
    done
    log_info "Including tags: $INCLUDE_TAGS"
fi

if [ -n "$EXCLUDE_TAGS" ]; then
    IFS=' ' read -ra TAGS <<< "$(echo "$EXCLUDE_TAGS" | tr ',' ' ')"
    for tag in "${TAGS[@]}"; do
        ROBOT_ARGS+=("--exclude" "$tag")
    done
    log_info "Excluding tags: $EXCLUDE_TAGS"
fi

# 添加變數檔案
if [ -n "$VARIABLE_FILE" ]; then
    IFS=',' read -ra VAR_FILES <<< "$VARIABLE_FILE"
    for var_file in "${VAR_FILES[@]}"; do
        var_file=$(echo "$var_file" | xargs)
        if [ -f "$var_file" ]; then
            ROBOT_ARGS+=("--variablefile" "$var_file")
            log_info "Using variable file: $var_file"
        else
            log_warning "Variable file not found: $var_file (skipping)"
        fi
    done
fi

# 添加並行執行
if [ -n "$PROCESSES" ] && [ "$PROCESSES" != "1" ]; then
    ROBOT_ARGS+=("--processes" "$PROCESSES")
    log_info "Running tests in parallel with $PROCESSES processes"
fi

# 添加超時
if [ -n "$TIMEOUT" ]; then
    ROBOT_ARGS+=("--testtimeout" "$TIMEOUT")
    log_info "Test timeout: $TIMEOUT"
fi

# 添加變數
for var in "${ROBOT_VARIABLES[@]}"; do
    ROBOT_ARGS+=("--variable" "$var")
done

# 添加測試目錄
ROBOT_ARGS+=("$TEST_DIR")

# 執行 robot 命令
log_info "Executing: robot ${ROBOT_ARGS[*]}"
robot "${ROBOT_ARGS[@]}" || EXIT_CODE=$?

# 輸出結果
if [ $EXIT_CODE -eq 0 ]; then
    log_success "All tests passed!"
else
    log_error "Tests failed with exit code: $EXIT_CODE"
fi

# 檢查報告檔案是否存在
if [ -f "$REPORT_DIR/report.html" ]; then
    log_success "Test report generated: $REPORT_DIR/report.html"
else
    log_warning "Test report not found at expected location: $REPORT_DIR/report.html"
fi

# 輸出報告路徑（供後續步驟使用）
# 使用新的 GITHUB_OUTPUT 方式（向後相容舊的 set-output）
if [ -n "${GITHUB_OUTPUT:-}" ]; then
    echo "exit_code=$EXIT_CODE" >> "$GITHUB_OUTPUT"
    echo "report_path=$REPORT_DIR" >> "$GITHUB_OUTPUT"
else
    # 向後相容舊的 set-output 方式
    echo "::set-output name=exit_code::$EXIT_CODE"
    echo "::set-output name=report_path::$REPORT_DIR"
fi

# 退出時使用相同的退出碼
exit $EXIT_CODE

