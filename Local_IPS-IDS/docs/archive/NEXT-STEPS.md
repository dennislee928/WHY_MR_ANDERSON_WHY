# 下一步執行計劃

> **當前階段**: 階段 2 - 建立完整前端結構  
> **完成度**: 60%

---

## 🎯 立即要做的事

### 1. 完成前端TailwindCSS配置

```bash
cd Application/Fe
npm install
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

**需要創建的檔案**:
- `tailwind.config.js`
- `postcss.config.js`

### 2. 創建缺少的前端組件

**Layout組件** (`Application/Fe/components/layout/MainLayout.tsx`):
- 導航欄
- 側邊欄
- 頁面容器

**Loading組件** (`Application/Fe/components/ui/loading.tsx`):
- 載入動畫
- 骨架屏

### 3. 完善API服務層

**創建** (`Application/Fe/services/api.ts`):
- 完整的API客戶端
- 錯誤處理
- 請求攔截器

### 4. 測試前端構建

```bash
cd Application/Fe
npm run build
npm run dev
```

---

## 📋 詳細執行命令

### Windows用戶

```powershell
# 進入前端目錄
cd Application\Fe

# 安裝依賴
npm install

# 安裝開發依賴
npm install -D tailwindcss postcss autoprefixer @types/node @types/react @types/react-dom typescript eslint eslint-config-next

# 初始化TailwindCSS
npx tailwindcss init -p

# 測試開發模式
npm run dev

# 測試構建
npm run build
```

### Linux/macOS用戶

```bash
# 進入前端目錄
cd Application/Fe

# 安裝依賴
npm install

# 安裝開發依賴
npm install -D tailwindcss postcss autoprefixer @types/node @types/react @types/react-dom typescript eslint eslint-config-next

# 初始化TailwindCSS
npx tailwindcss init -p

# 測試開發模式
npm run dev

# 測試構建
npm run build
```

---

## 🔧 需要手動創建的檔案

### 1. Tailwind配置

**`Application/Fe/tailwind.config.js`**:
```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        'pandora': {
          50: '#f5f7ff',
          100: '#ebf0fe',
          200: '#d6e1fd',
          300: '#b3c9fb',
          400: '#8aa5f8',
          500: '#667eea',
          600: '#5568d3',
          700: '#4553b8',
          800: '#3a4595',
          900: '#333d7a',
        },
      },
    },
  },
  plugins: [],
}
```

### 2. PostCSS配置

**`Application/Fe/postcss.config.js`**:
```javascript
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

### 3. Layout組件

**`Application/Fe/components/layout/MainLayout.tsx`**:
```typescript
import React from 'react'

interface MainLayoutProps {
  children: React.ReactNode
}

export default function MainLayout({ children }: MainLayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4">
          <h1 className="text-2xl font-bold text-gray-900">
            Pandora Box Console IDS-IPS
          </h1>
        </div>
      </header>
      <main className="max-w-7xl mx-auto px-4 py-6">
        {children}
      </main>
    </div>
  )
}
```

---

## ✅ 驗證清單

完成後，請確認以下項目：

- [ ] `npm install` 成功執行
- [ ] TailwindCSS 配置正確
- [ ] `npm run dev` 可以啟動開發伺服器
- [ ] 瀏覽器可以訪問 http://localhost:3001
- [ ] Dashboard 頁面正常顯示
- [ ] 沒有TypeScript錯誤
- [ ] 沒有ESLint錯誤
- [ ] `npm run build` 成功構建
- [ ] 構建產物在 `.next/` 目錄

---

## 🚨 可能遇到的問題

### 問題 1: npm install 失敗

**解決方案**:
```bash
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### 問題 2: TypeScript錯誤

**解決方案**:
```bash
npm install -D typescript @types/node @types/react @types/react-dom
```

### 問題 3: 端口被占用

**解決方案**:
```bash
# Windows
netstat -ano | findstr :3001
taskkill /PID <PID> /F

# Linux/macOS
lsof -ti:3001 | xargs kill -9
```

### 問題 4: 找不到模組

**解決方案**:
檢查 `tsconfig.json` 中的路徑配置是否正確

---

## 📞 需要幫助？

如果遇到問題，請：

1. 檢查 [PROJECT-RESTRUCTURE-PROGRESS.md](PROJECT-RESTRUCTURE-PROGRESS.md)
2. 查看錯誤訊息
3. 提供完整的錯誤日誌

---

**更新時間**: 2025-10-09 09:40

