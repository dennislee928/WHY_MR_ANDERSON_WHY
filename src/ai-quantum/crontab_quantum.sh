#!/bin/bash
# Pandora Quantum Analysis Cron Jobs
# 添加到 crontab: crontab -e

# ========== 每日量子分析 (每天凌晨 2:00) ==========
# 0 2 * * * cd /app && python scheduled_quantum_analysis.py daily >> /app/logs/daily_analysis.log 2>&1

# ========== 每週量子訓練 (每週日凌晨 3:00) ==========
# 0 3 * * 0 cd /app && python scheduled_quantum_analysis.py weekly >> /app/logs/weekly_training.log 2>&1

# ========== 每月批次分析 (每月1號凌晨 4:00) ==========
# 0 4 1 * * cd /app && python scheduled_quantum_analysis.py monthly >> /app/logs/monthly_analysis.log 2>&1

# ========== Crontab 語法說明 ==========
# 格式: 分鐘 小時 日 月 星期 命令
# *: 任何值
# */5: 每5個單位
# 1-5: 範圍
# 1,3,5: 列表

# ========== 範例使用 ==========
echo "Pandora Quantum Analysis Cron Jobs"
echo ""
echo "安裝步驟:"
echo "1. chmod +x crontab_quantum.sh"
echo "2. 編輯此文件，取消註釋需要的作業"
echo "3. 複製 cron 表達式到 crontab"
echo "4. crontab -e"
echo ""
echo "可用的分析類型:"
echo "  - daily: 每日高風險事件重新評估"
echo "  - weekly: 每週模型重訓練"
echo "  - monthly: 每月批次深度分析"

