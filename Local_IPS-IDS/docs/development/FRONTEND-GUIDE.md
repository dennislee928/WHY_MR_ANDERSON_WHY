# 前端開發指南

> **技術棧**: Next.js 14 + TypeScript + TailwindCSS

---

## 🚀 快速開始

```bash
cd Application/Fe
npm install
npm run dev
```

---

## 📁 目錄結構

詳見 [Application/Fe/README.md](../../Application/Fe/README.md)

---

## 🎨 組件開發

### 創建新組件

```typescript
// components/security/EventList.tsx
import React from 'react'
import { Card } from '../ui/card'

export function EventList() {
  return (
    <Card>
      <h2>安全事件列表</h2>
    </Card>
  )
}
```

---

**維護**: Development Team

