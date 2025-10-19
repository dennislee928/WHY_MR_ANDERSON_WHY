# 🚀 從這裡開始！

> **第一次使用 Application/?** 看這個就對了！

---

## ⚡ 10 秒決定：我要...

### 選項 A: 快速體驗（最簡單）

```bash
cd Application
./docker-start.sh    # 或 .\docker-start.ps1 (Windows)
```

✅ **2分鐘啟動完整系統**  
✅ **11個服務自動運行**  
✅ **訪問 http://localhost:3001**

**需要**: Docker

---

### 選項 B: 本地開發（完整）

```bash
cd Application

# Windows
.\build-local.ps1

# Linux/macOS  
./build-local.sh

cd dist
./start.sh  # 或 start.bat
```

✅ **編譯原生程式**  
✅ **最佳效能**  
✅ **適合開發和自訂**

**需要**: Go 1.24+, Node.js 18+

---

### 選項 C: 只開發後端

```bash
cd Application/be
make all
make run-agent
```

✅ **最快編譯**  
✅ **專注 Go 開發**

**需要**: Go 1.24+

---

## 📊 對照表

| 項目 | Docker | 本地構建 | 後端專用 |
|------|--------|----------|----------|
| **啟動時間** | 2分鐘 | 5分鐘 | 1分鐘 |
| **服務數量** | 11個 | 3個 | 1-3個 |
| **需要** | Docker | Go+Node | Go |
| **適合** | 測試 | 開發 | Go開發 |

---

## 🎯 我的建議

### 第一次使用？
👉 **選擇 Docker Compose（選項A）**  
最快看到完整系統運行

### 要開發程式？
👉 **選擇本地構建（選項B）**  
完整的開發體驗

### 只改後端？
👉 **選擇後端專用（選項C）**  
最快的編譯循環

---

## 📚 詳細說明

完整指南請看：
- [Application/README.md](README.md)
- [Application/DEPLOYMENT-OPTIONS.md](DEPLOYMENT-OPTIONS.md)
- [Application/DOCKER-ARCHITECTURE.md](DOCKER-ARCHITECTURE.md)

---

**選好了？立即開始！** 🚀

