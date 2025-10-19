# Pandora Box Console - 前端應用程式

完整的 Next.js 14 前端應用程式，用於地端部署。

## ✅ 完成項目

### 已實作功能

- ✅ **完整目錄結構**
  - components/（UI組件、Layout、Dashboard、Security）
  - pages/（Next.js頁面）
  - services/（API服務層）
  - hooks/（自定義Hooks）
  - types/（TypeScript類型定義）
  - lib/（工具函數）
  - styles/（全域樣式）

- ✅ **UI組件系統**
  - Card, CardHeader, CardTitle, CardContent
  - Button（多種變體）
  - Badge
  - Loading（載入動畫和骨架屏）
  - Alert（多種類型的提示）

- ✅ **佈局系統**
  - MainLayout（響應式導航、Header、Footer）
  - 支援移動端側邊欄

- ✅ **API服務層**
  - 統一的API客戶端
  - 錯誤處理
  - 完整的API方法封裝
  - WebSocket支援

- ✅ **自定義Hooks**
  - useSystemStatus（系統狀態管理）
  - useWebSocket（WebSocket連接管理）

- ✅ **工具函數庫**
  - 格式化函數（bytes, number, uptime, timestamp）
  - 樣式工具（classNames, colors）
  - 性能優化（debounce, throttle）

- ✅ **配置完整**
  - TailwindCSS（含自定義主題）
  - PostCSS
  - TypeScript
  - ESLint
  - Next.js
  - 環境變數

## 📁 目錄結構

```
Fe/
├── components/
│   ├── ui/                    # ✅ 基礎UI組件
│   │   ├── card.tsx
│   │   ├── button.tsx
│   │   ├── badge.tsx
│   │   ├── loading.tsx
│   │   └── alert.tsx
│   ├── layout/                # ✅ 佈局組件
│   │   └── MainLayout.tsx
│   ├── dashboard/             # ✅ 儀表板組件
│   │   └── Dashboard.tsx
│   └── security/              # 安全相關組件（待實作）
│
├── pages/                     # ✅ Next.js頁面
│   ├── _app.tsx
│   ├── _document.tsx
│   └── index.tsx
│
├── services/                  # ✅ API服務
│   └── api.ts
│
├── hooks/                     # ✅ 自定義Hooks
│   ├── useSystemStatus.ts
│   └── useWebSocket.ts
│
├── types/                     # ✅ TypeScript類型
│   └── index.ts
│
├── lib/                       # ✅ 工具函數
│   └── utils.ts
│
├── styles/                    # ✅ 樣式
│   └── globals.css
│
├── public/                    # ✅ 靜態資源
│   ├── favicon.ico
│   └── index.html
│
├── legacy/                    # 舊版程式碼存檔
│   └── web-old/
│
├── package.json               # ✅ NPM配置
├── tsconfig.json              # ✅ TypeScript配置
├── next.config.js             # ✅ Next.js配置
├── tailwind.config.js         # ✅ TailwindCSS配置
├── postcss.config.js          # ✅ PostCSS配置
├── .eslintrc.json             # ✅ ESLint配置
├── .gitignore                 # ✅ Git忽略規則
├── .env.example               # ✅ 環境變數範例
└── README.md                  # 本檔案
```

## 🚀 快速開始

### 1. 安裝依賴

```bash
npm install
```

### 2. 設定環境變數

```bash
cp .env.example .env.local
# 編輯 .env.local 設定API URL等
```

### 3. 開發模式

```bash
npm run dev
```

訪問 http://localhost:3001

### 4. 生產構建

```bash
npm run build
npm run start
```

## 📦 主要依賴

| 套件 | 版本 | 說明 |
|------|------|------|
| next | ^14.1.0 | React框架 |
| react | ^18.2.0 | UI庫 |
| typescript | ^5.3.3 | 類型系統 |
| tailwindcss | ^3.4.1 | CSS框架 |
| lucide-react | ^0.312.0 | 圖示庫 |

## 🔧 可用腳本

| 命令 | 說明 |
|------|------|
| `npm run dev` | 啟動開發伺服器（3001端口） |
| `npm run build` | 構建生產版本 |
| `npm run start` | 啟動生產伺服器 |
| `npm run lint` | 執行ESLint檢查 |
| `npm run type-check` | 執行TypeScript類型檢查 |
| `npm run format` | 格式化程式碼（如已安裝Prettier） |

## 🎨 主題配置

TailwindCSS 自定義主題色：

```javascript
colors: {
  'pandora': {
    50: '#f5f7ff',
    500: '#667eea',
    900: '#333d7a',
  }
}
```

## 📝 API使用範例

### 使用API服務

```typescript
import { fetchSystemStatus, blockNetwork } from '../services/api'

// 獲取系統狀態
const status = await fetchSystemStatus()

// 阻斷網路
await blockNetwork()
```

### 使用Hooks

```typescript
import { useSystemStatus } from '../hooks/useSystemStatus'

function Component() {
  const { status, loading, error } = useSystemStatus(30000)
  
  if (loading) return <Loading />
  if (error) return <Alert type="error" message={error} />
  
  return <div>{status.agent.status}</div>
}
```

## 🐛 已知問題

無重大問題。

## 📚 下一步

1. 安裝依賴：`cd Application/Fe && npm install`
2. 測試開發模式：`npm run dev`
3. 測試構建：`npm run build`
4. 實作更多頁面（安全事件、設定等）

## 🤝 開發指南

### 新增頁面

在 `pages/` 目錄創建新檔案：

```typescript
// pages/security.tsx
import MainLayout from '../components/layout/MainLayout'

export default function Security() {
  return (
    <MainLayout>
      <h1>安全事件</h1>
    </MainLayout>
  )
}
```

### 新增組件

在 `components/` 對應目錄創建：

```typescript
// components/security/EventList.tsx
export function EventList() {
  return <div>事件列表</div>
}
```

---

**狀態**: ✅ 階段2完成（100%）  
**最後更新**: 2025-10-09 09:45
