# Pandora Box Console v3.3.0 部署檢查清單

## ✅ 部署前檢查

### 1. 環境準備
- [ ] Docker 20.10+ 已安裝
- [ ] Docker Compose 2.0+ 已安裝
- [ ] 至少 8GB RAM 可用
- [ ] 至少 20GB 磁碟空間
- [ ] 網路連接正常

### 2. 配置文件
- [ ] `configs/ui-config.yaml` 存在
- [ ] `configs/agent-config.yaml` 存在
- [ ] Nginx 配置已更新（`axiom-ui` → `axiom-be`）
- [ ] IBM Quantum Token 已設置（可選）

---

## 🚀 部署步驟

### 步驟 1: 清理舊容器
```bash
cd Application

# 停止所有服務
docker-compose down

# 刪除舊的 axiom-ui 容器（如果存在）
docker rm -f axiom-ui

# 清理未使用的映像（可選）
docker system prune -f
```

### 步驟 2: 構建新映像
```bash
# 構建 axiom-be（獨立後端）
docker-compose build --no-cache axiom-be

# 構建 cyber-ai-quantum（量子服務）
docker-compose build --no-cache cyber-ai-quantum

# 或構建所有服務
docker-compose build --no-cache
```

### 步驟 3: 啟動服務
```bash
# 啟動所有服務
docker-compose up -d

# 或使用腳本
./docker-start.sh
```

### 步驟 4: 驗證部署
```bash
# 檢查所有容器狀態
docker-compose ps

# 應該看到 14 個容器運行中：
# ✅ portainer
# ✅ axiom-be
# ✅ pandora-agent
# ✅ cyber-ai-quantum
# ✅ prometheus
# ✅ grafana
# ✅ loki
# ✅ alertmanager
# ✅ rabbitmq
# ✅ postgres
# ✅ redis
# ✅ node-exporter
# ✅ promtail
# ✅ nginx
```

---

## ✅ 部署後驗證

### 1. Portainer 訪問
```bash
# 訪問 Portainer
curl http://localhost:9000/api/status

# 或瀏覽器打開
http://localhost:9000
```

**預期結果**: Portainer 登入頁面

### 2. Axiom BE API 測試
```bash
# 健康檢查
curl http://localhost:3001/api/v1/health

# 系統狀態
curl http://localhost:3001/api/v1/status

# Swagger 文檔
curl http://localhost:3001/swagger.json
```

**預期結果**: 所有端點返回 200 OK

### 3. Cyber AI/Quantum 測試
```bash
# 健康檢查
curl http://localhost:8000/health

# 系統狀態
curl http://localhost:8000/api/v1/status

# FastAPI 文檔
curl http://localhost:8000/docs
```

**預期結果**: 所有端點正常

### 4. IBM Quantum 連接測試（可選）
```bash
cd Experimental/cyber-ai-quantum

# 設置 Token
export IBM_QUANTUM_TOKEN=7PzS0AdaFB83OVn7WghkdLupzthCe1fTl3yI4Qrp7G6o

# 測試連接
python test_ibm_connection.py
```

**預期結果**: 顯示可用的量子後端列表

### 5. 監控服務測試
```bash
# Prometheus
curl http://localhost:9090/-/healthy

# Grafana
curl http://localhost:3000/api/health

# Loki
curl http://localhost:3100/ready

# AlertManager
curl http://localhost:9093/-/healthy
```

**預期結果**: 所有服務返回 healthy

---

## 🔧 常見問題修復

### 問題 1: 端口衝突（Port already allocated）

**症狀**: 
```
Bind for 0.0.0.0:3001 failed: port is already allocated
```

**解決方案**:
```bash
# 找到佔用端口的容器
docker ps | grep 3001

# 停止舊容器
docker stop axiom-ui
docker rm axiom-ui

# 重新啟動
docker-compose up -d axiom-be
```

### 問題 2: Nginx 找不到上游服務

**症狀**:
```
host not found in upstream "axiom-ui:3001"
```

**解決方案**:
```bash
# 已修復：更新 nginx 配置
# configs/nginx/nginx.conf: axiom-ui:3001 → axiom-be:3001
# configs/nginx/default-paas.conf: axiom-ui:3001 → axiom-be:3001

# 重啟 nginx
docker-compose restart nginx
```

### 問題 3: 容器無法啟動

**解決方案**:
```bash
# 查看日誌
docker-compose logs [service_name]

# 重新構建
docker-compose build --no-cache [service_name]

# 重新啟動
docker-compose up -d [service_name]
```

---

## 📊 健康檢查清單

### 容器健康狀態
- [ ] portainer: healthy
- [ ] axiom-be: healthy
- [ ] pandora-agent: healthy
- [ ] cyber-ai-quantum: healthy
- [ ] prometheus: healthy
- [ ] grafana: healthy
- [ ] loki: healthy
- [ ] alertmanager: healthy
- [ ] rabbitmq: healthy
- [ ] postgres: healthy
- [ ] redis: healthy
- [ ] node-exporter: up
- [ ] promtail: healthy
- [ ] nginx: healthy

### API 端點測試
- [ ] http://localhost:9000 (Portainer)
- [ ] http://localhost:3001/api/v1/health (Axiom BE)
- [ ] http://localhost:3001/swagger (Swagger UI)
- [ ] http://localhost:8000/health (AI/Quantum)
- [ ] http://localhost:8000/docs (FastAPI Docs)
- [ ] http://localhost:3000 (Grafana)
- [ ] http://localhost:9090 (Prometheus)
- [ ] http://localhost:15672 (RabbitMQ)

### 數據庫連接
- [ ] PostgreSQL: localhost:5432 (pandora/pandora123)
- [ ] Redis: localhost:6379 (密碼: pandora123)

---

## 🎯 Portainer 初次設置

### 步驟 1: 訪問 Portainer
```
瀏覽器打開: http://localhost:9000
```

### 步驟 2: 創建管理員帳號
```
用戶名: admin
密碼: pandora_portainer_2025! (至少 12 字元)
```

### 步驟 3: 連接環境
```
選擇: Get Started
環境: Local (自動檢測)
```

### 步驟 4: 驗證
```
Dashboard → 應該看到:
- Stacks: 1 (application)
- Containers: 14 (14 running)
- Images: 12+
- Volumes: 14
```

---

## 📈 監控設置

### 1. Grafana 初次登入
```
URL: http://localhost:3000
用戶名: admin
密碼: pandora123
```

### 2. 配置數據源（已自動配置）
- ✅ Prometheus
- ✅ Loki
- ⚠️ AlertManager（需手動配置）

### 3. 導入儀表板
```
Configuration → Data Sources → 驗證所有數據源
Dashboards → Browse → 查看預設儀表板
```

---

## 🔒 安全檢查

### 1. 更改預設密碼
```bash
# Grafana
訪問 http://localhost:3000 → Profile → Change Password

# Portainer
訪問 http://localhost:9000 → Account → Change Password

# PostgreSQL（在生產環境）
docker exec -it postgres psql -U pandora -c "ALTER USER pandora PASSWORD 'new_password';"
```

### 2. 配置防火牆（生產環境）
```bash
# 僅允許必要端口
sudo ufw allow 3000  # Grafana
sudo ufw allow 3001  # Axiom BE
sudo ufw allow 8000  # AI/Quantum
sudo ufw allow 9000  # Portainer
sudo ufw enable
```

### 3. 啟用 HTTPS
```bash
# 使用 Nginx HTTPS（已配置）
# 訪問: https://localhost:443

# Portainer HTTPS
# 訪問: https://localhost:9443
```

---

## 📝 維護任務

### 每日
- [ ] 訪問 Portainer Dashboard
- [ ] 檢查所有容器運行中
- [ ] 查看資源使用率

### 每週
- [ ] 查看 Grafana 指標趨勢
- [ ] 檢查告警歷史
- [ ] 清理未使用的 Docker 映像

### 每月
- [ ] 備份 PostgreSQL 數據
- [ ] 備份 Grafana 配置
- [ ] 更新 Docker 映像
- [ ] 檢查磁碟空間

---

## 🎉 部署成功標誌

當您看到以下情況時，部署成功：

✅ **Portainer Dashboard**:
- 14 個容器全部顯示為綠色（運行中）
- 無錯誤或警告

✅ **Axiom BE**:
- http://localhost:3001/api/v1/health 返回 200
- http://localhost:3001/swagger 顯示 Swagger UI

✅ **Cyber AI/Quantum**:
- http://localhost:8000/health 返回 200
- http://localhost:8000/docs 顯示 FastAPI 文檔

✅ **Grafana**:
- http://localhost:3000 可訪問
- 所有數據源顯示綠色

✅ **Nginx**:
- 無 "host not found" 錯誤
- 容器狀態為 "healthy"

---

**🎊 恭喜！Pandora Box Console v3.3.0 部署成功！** 🎊

---

**維護者**: Pandora Security Team  
**版本**: v3.3.0  
**最後更新**: 2025-01-14

