# Pandora Box Console - 完整專案重構腳本
# 此腳本會將專案重組為地端部署結構

$ErrorActionPreference = "Stop"

Write-Host "=================================" -ForegroundColor Cyan
Write-Host "  Pandora Box Console 專案重構  " -ForegroundColor Cyan
Write-Host "=================================" -ForegroundColor Cyan
Write-Host ""

$ROOT_DIR = (Get-Location).Path

# 步驟 1-3: 已完成
Write-Host "✓ 步驟 1-3: 已完成目錄結構和基礎檔案" -ForegroundColor Green

# 步驟 4: 創建後端引用結構
Write-Host "" 
Write-Host "步驟 4: 設定後端引用結構..." -ForegroundColor Yellow

# 在 Application/be/ 中創建 go.mod 引用根目錄
$goModContent = @"
// 此 go.mod 引用根目錄的模組
// 編譯時使用根目錄的 go.mod

module github.com/your-org/pandora_box_console_IDS-IPS/Application/be

go 1.24

// 引用根目錄模組
replace github.com/your-org/pandora_box_console_IDS-IPS => ../..

require github.com/your-org/pandora_box_console_IDS-IPS v0.0.0
"@

$goModContent | Out-File -FilePath "Application\be\go.mod" -Encoding UTF8

Write-Host "✓ 已創建 Application/be/go.mod" -ForegroundColor Green

# 步驟 5: 創建完整的前端檔案結構
Write-Host ""
Write-Host "步驟 5: 創建完整前端框架..." -ForegroundColor Yellow

# 創建 Next.js 頁面
$indexPage = @"
import React from 'react'
import Dashboard from '@/components/dashboard/Dashboard'

export default function Home() {
  const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'
  
  return <Dashboard apiBaseUrl={apiBaseUrl} />
}
"@

$indexPage | Out-File -FilePath "Application\Fe\pages\index.tsx" -Encoding UTF8

# 創建 _app.tsx
$appPage = @"
import '@/styles/globals.css'
import type { AppProps } from 'next/app'

export default function App({ Component, pageProps }: AppProps) {
  return <Component {...pageProps} />
}
"@

$appPage | Out-File -FilePath "Application\Fe\pages\_app.tsx" -Encoding UTF8

# 創建全域 CSS
$globalCSS = @"
@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  --pandora-600: #667eea;
  --pandora-500: #764ba2;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
"@

$globalCSS | Out-File -FilePath "Application\Fe\styles\globals.css" -Encoding UTF8

Write-Host "✓ 已創建前端頁面和樣式" -ForegroundColor Green

# 步驟 6: 創建前端 UI 組件
Write-Host ""
Write-Host "步驟 6: 創建 UI 組件..." -ForegroundColor Yellow

# Card 組件
$cardComponent = @"
import React from 'react'

interface CardProps {
  children: React.ReactNode
  className?: string
}

export function Card({ children, className = '' }: CardProps) {
  return (
    <div className={\`bg-white rounded-lg shadow-md \${className}\`}>
      {children}
    </div>
  )
}

export function CardHeader({ children, className = '' }: CardProps) {
  return <div className={\`p-6 pb-4 \${className}\`}>{children}</div>
}

export function CardTitle({ children, className = '' }: CardProps) {
  return <h3 className={\`text-lg font-semibold \${className}\`}>{children}</h3>
}

export function CardContent({ children, className = '' }: CardProps) {
  return <div className={\`p-6 pt-0 \${className}\`}>{children}</div>
}
"@

$cardComponent | Out-File -FilePath "Application\Fe\components\ui\card.tsx" -Encoding UTF8

# Button 組件
$buttonComponent = @"
import React from 'react'

interface ButtonProps {
  children: React.ReactNode
  onClick?: () => void
  className?: string
  variant?: 'default' | 'outline' | 'destructive'
  disabled?: boolean
}

export function Button({ 
  children, 
  onClick, 
  className = '', 
  variant = 'default',
  disabled = false 
}: ButtonProps) {
  const baseClasses = 'px-4 py-2 rounded-md font-medium transition-colors'
  const variantClasses = {
    default: 'bg-blue-600 text-white hover:bg-blue-700',
    outline: 'border border-gray-300 text-gray-700 hover:bg-gray-50',
    destructive: 'bg-red-600 text-white hover:bg-red-700'
  }
  
  return (
    <button
      onClick={onClick}
      disabled={disabled}
      className={\`\${baseClasses} \${variantClasses[variant]} \${className} \${disabled ? 'opacity-50 cursor-not-allowed' : ''}\`}
    >
      {children}
    </button>
  )
}
"@

$buttonComponent | Out-File -FilePath "Application\Fe\components\ui\button.tsx" -Encoding UTF8

# Badge 組件
$badgeComponent = @"
import React from 'react'

interface BadgeProps {
  children: React.ReactNode
  className?: string
}

export function Badge({ children, className = '' }: BadgeProps) {
  return (
    <span className={\`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium \${className}\`}>
      {children}
    </span>
  )
}
"@

$badgeComponent | Out-File -FilePath "Application\Fe\components\ui\badge.tsx" -Encoding UTF8

Write-Host "✓ 已創建 UI 組件" -ForegroundColor Green

# 步驟 7: 創建服務層
Write-Host ""
Write-Host "步驟 7: 創建服務層..." -ForegroundColor Yellow

$apiService = @"
const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

export async function fetchSystemStatus() {
  const response = await fetch(\`\${API_BASE_URL}/v1/status\`)
  if (!response.ok) {
    throw new Error('Failed to fetch system status')
  }
  return response.json()
}

export async function controlNetwork(action: 'block' | 'unblock') {
  const response = await fetch(\`\${API_BASE_URL}/v1/control/network\`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ action }),
  })
  if (!response.ok) {
    throw new Error(\`Failed to \${action} network\`)
  }
  return response.json()
}

export async function fetchEvents(limit = 10) {
  const response = await fetch(\`\${API_BASE_URL}/v1/events?limit=\${limit}\`)
  if (!response.ok) {
    throw new Error('Failed to fetch events')
  }
  return response.json()
}
"@

$apiService | Out-File -FilePath "Application\Fe\services\api.ts" -Encoding UTF8

Write-Host "✓ 已創建服務層" -ForegroundColor Green

# 步驟 8: 創建 TypeScript 類型定義
Write-Host ""
Write-Host "步驟 8: 創建類型定義..." -ForegroundColor Yellow

$typesFile = @"
export interface SystemStatus {
  agent: {
    status: 'online' | 'offline' | 'error'
    lastSeen: string
    uptime: string
  }
  console: {
    status: 'online' | 'offline' | 'error'
    requests: number
    responseTime: number
  }
  network: {
    blocked: boolean
    blockTime: string
    unlockTime: string
  }
  metrics: {
    totalAlerts: number
    criticalAlerts: number
    systemLoad: number
    memoryUsage: number
  }
}

export interface SecurityEvent {
  id: string
  timestamp: number
  type: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  message: string
  source?: string
  destination?: string
}

export interface NetworkControlRequest {
  action: 'block' | 'unblock'
}

export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: string
  message?: string
}
"@

$typesFile | Out-File -FilePath "Application\Fe\types\index.ts" -Encoding UTF8

Write-Host "✓ 已創建類型定義" -ForegroundColor Green

# 步驟 9: 更新 Dashboard 組件以使用新結構
Write-Host ""
Write-Host "步驟 9: 更新 Dashboard 組件..." -ForegroundColor Yellow

$updatedDashboard = Get-Content "Application\Fe\components\dashboard\Dashboard.tsx" -Raw
$updatedDashboard = $updatedDashboard -replace "@/components/ui/card", "../ui/card"
$updatedDashboard = $updatedDashboard -replace "@/components/ui/button", "../ui/button"
$updatedDashboard = $updatedDashboard -replace "@/components/ui/badge", "../ui/badge"
$updatedDashboard | Out-File -FilePath "Application\Fe\components\dashboard\Dashboard.tsx" -Encoding UTF8

Write-Host "✓ 已更新 Dashboard 組件" -ForegroundColor Green

# 步驟 10: 創建環境變數範例檔案
Write-Host ""
Write-Host "步驟 10: 創建配置檔案..." -ForegroundColor Yellow

$envExample = @"
# API Configuration
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws

# Monitoring URLs
NEXT_PUBLIC_GRAFANA_URL=http://localhost:3000
NEXT_PUBLIC_PROMETHEUS_URL=http://localhost:9090
NEXT_PUBLIC_LOKI_URL=http://localhost:3100

# App Configuration
NEXT_PUBLIC_APP_NAME=Pandora Box Console IDS-IPS
NEXT_PUBLIC_APP_VERSION=3.0.0
"@

$envExamplePath = "Application\Fe\.env.example"
$envLocalPath = "Application\Fe\.env.local"
$envExample | Out-File -FilePath $envExamplePath -Encoding UTF8
$envExample | Out-File -FilePath $envLocalPath -Encoding UTF8

Write-Host "✓ 已創建環境變數檔案" -ForegroundColor Green

# 完成
Write-Host ""
Write-Host "=================================" -ForegroundColor Green
Write-Host "  ✅ 專案重構完成！" -ForegroundColor Green
Write-Host "=================================" -ForegroundColor Green
Write-Host ""
Write-Host "新的專案結構：" -ForegroundColor Cyan
Write-Host "  Application/Fe/  - 完整的 Next.js 前端" -ForegroundColor White
Write-Host "  Application/be/  - 後端引用結構" -ForegroundColor White
Write-Host ""
Write-Host "下一步：" -ForegroundColor Yellow
Write-Host "  1. cd Application\Fe; npm install" -ForegroundColor White
Write-Host "  2. npm run dev  # 啟動前端開發伺服器" -ForegroundColor White
Write-Host "  3. cd Application\be; make all  # 編譯後端" -ForegroundColor White
Write-Host ""

