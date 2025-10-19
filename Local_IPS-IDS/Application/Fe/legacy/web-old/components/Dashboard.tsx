import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { 
  Activity, 
  Shield, 
  AlertTriangle, 
  CheckCircle, 
  Clock,
  Wifi,
  WifiOff
} from 'lucide-react'

interface DashboardProps {
  apiBaseUrl: string
}

interface SystemStatus {
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

export default function Dashboard({ apiBaseUrl }: DashboardProps) {
  const [systemStatus, setSystemStatus] = useState<SystemStatus | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchSystemStatus = async () => {
      try {
        setLoading(true)
        const response = await fetch(`${apiBaseUrl}/v1/status`)
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        
        const data = await response.json()
        setSystemStatus(data)
        setError(null)
      } catch (err) {
        setError(err instanceof Error ? err.message : '未知錯誤')
        console.error('獲取系統狀態失敗:', err)
      } finally {
        setLoading(false)
      }
    }

    // 初始載入
    fetchSystemStatus()

    // 每 30 秒更新一次
    const interval = setInterval(fetchSystemStatus, 30000)

    return () => clearInterval(interval)
  }, [apiBaseUrl])

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online':
        return <CheckCircle className="h-5 w-5 text-green-500" />
      case 'offline':
        return <WifiOff className="h-5 w-5 text-red-500" />
      case 'error':
        return <AlertTriangle className="h-5 w-5 text-yellow-500" />
      default:
        return <Clock className="h-5 w-5 text-gray-500" />
    }
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'online':
        return <Badge className="bg-green-100 text-green-800">線上</Badge>
      case 'offline':
        return <Badge className="bg-red-100 text-red-800">離線</Badge>
      case 'error':
        return <Badge className="bg-yellow-100 text-yellow-800">錯誤</Badge>
      default:
        return <Badge className="bg-gray-100 text-gray-800">未知</Badge>
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-pandora-600"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Card className="w-96">
          <CardContent className="pt-6">
            <div className="text-center">
              <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-900 mb-2">載入失敗</h3>
              <p className="text-gray-600 mb-4">{error}</p>
              <Button onClick={() => window.location.reload()}>
                重新載入
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        {/* 頁面標題 */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Pandora Box Console
          </h1>
          <p className="text-gray-600">
            IDS/IPS 網路安全監控系統
          </p>
        </div>

        {/* 系統狀態概覽 */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {/* Agent 狀態 */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Agent 狀態</CardTitle>
              {systemStatus && getStatusIcon(systemStatus.agent.status)}
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between">
                <div>
                  <div className="text-2xl font-bold">
                    {systemStatus?.agent.status === 'online' ? '正常' : '異常'}
                  </div>
                  <p className="text-xs text-gray-600">
                    最後更新: {systemStatus?.agent.lastSeen}
                  </p>
                </div>
                {systemStatus && getStatusBadge(systemStatus.agent.status)}
              </div>
            </CardContent>
          </Card>

          {/* Console 狀態 */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Console 狀態</CardTitle>
              <Activity className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {systemStatus?.console.requests || 0}
              </div>
              <p className="text-xs text-gray-600">
                響應時間: {systemStatus?.console.responseTime || 0}ms
              </p>
            </CardContent>
          </Card>

          {/* 網路狀態 */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">網路狀態</CardTitle>
              {systemStatus?.network.blocked ? (
                <WifiOff className="h-4 w-4 text-red-500" />
              ) : (
                <Wifi className="h-4 w-4 text-green-500" />
              )}
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {systemStatus?.network.blocked ? '已阻斷' : '正常'}
              </div>
              <p className="text-xs text-gray-600">
                阻斷時間: {systemStatus?.network.blockTime}
              </p>
            </CardContent>
          </Card>

          {/* 告警統計 */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">告警統計</CardTitle>
              <Shield className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">
                {systemStatus?.metrics.totalAlerts || 0}
              </div>
              <p className="text-xs text-gray-600">
                嚴重告警: {systemStatus?.metrics.criticalAlerts || 0}
              </p>
            </CardContent>
          </Card>
        </div>

        {/* 詳細資訊 */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* 系統指標 */}
          <Card>
            <CardHeader>
              <CardTitle>系統指標</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium">系統負載</span>
                  <span className="text-sm text-gray-600">
                    {systemStatus?.metrics.systemLoad || 0}%
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-pandora-600 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${systemStatus?.metrics.systemLoad || 0}%` }}
                  ></div>
                </div>
                
                <div className="flex justify-between items-center">
                  <span className="text-sm font-medium">記憶體使用</span>
                  <span className="text-sm text-gray-600">
                    {systemStatus?.metrics.memoryUsage || 0}%
                  </span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-pandora-500 h-2 rounded-full transition-all duration-300"
                    style={{ width: `${systemStatus?.metrics.memoryUsage || 0}%` }}
                  ></div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* 快速操作 */}
          <Card>
            <CardHeader>
              <CardTitle>快速操作</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                <Button className="w-full" variant="outline">
                  重新載入系統狀態
                </Button>
                <Button className="w-full" variant="outline">
                  查看詳細日誌
                </Button>
                <Button className="w-full" variant="outline">
                  系統設定
                </Button>
                <Button className="w-full" variant="destructive">
                  緊急停止
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
