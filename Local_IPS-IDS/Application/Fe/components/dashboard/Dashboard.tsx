import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card'
import { Button } from '../ui/button'
import { Badge } from '../ui/badge'
import { 
  Activity, 
  Shield, 
  AlertTriangle, 
  CheckCircle, 
  Clock,
  Wifi,
  WifiOff,
  Server,
  Database,
  Cpu,
  HardDrive,
  Network,
  Users,
  Eye,
  Lock,
  Unlock,
  RefreshCw,
  Settings,
  BarChart3,
  TrendingUp,
  TrendingDown,
  Zap
} from 'lucide-react'

interface DashboardProps {
  apiBaseUrl: string
}

interface SystemStatus {
  agent: {
    status: 'online' | 'offline' | 'error'
    lastSeen: string
    uptime: string
    version: string
    cpuUsage: number
    memoryUsage: number
    diskUsage: number
  }
  console: {
    status: 'online' | 'offline' | 'error'
    requests: number
    responseTime: number
    activeConnections: number
    errorRate: number
  }
  network: {
    blocked: boolean
    blockTime: string
    unlockTime: string
    totalTraffic: number
    blockedIPs: number
    allowedIPs: number
  }
  security: {
    totalAlerts: number
    criticalAlerts: number
    warningAlerts: number
    infoAlerts: number
    threatsBlocked: number
    lastThreat: string
  }
  monitoring: {
    prometheus: boolean
    grafana: boolean
    loki: boolean
    alertmanager: boolean
  }
  devices: {
    total: number
    online: number
    offline: number
    lastActivity: string
  }
}

interface Alert {
  id: string
  level: 'critical' | 'warning' | 'info'
  message: string
  timestamp: string
  source: string
  resolved: boolean
}

export default function Dashboard({ apiBaseUrl }: DashboardProps) {
  const [systemStatus, setSystemStatus] = useState<SystemStatus | null>(null)
  const [alerts, setAlerts] = useState<Alert[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [activeTab, setActiveTab] = useState<'overview' | 'security' | 'network' | 'monitoring'>('overview')

  useEffect(() => {
    const fetchSystemStatus = async () => {
      try {
        setLoading(true)
        const response = await fetch(`${apiBaseUrl}/api/v1/status`)
        
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

    const fetchAlerts = async () => {
      try {
        const response = await fetch(`${apiBaseUrl}/api/v1/alerts`)
        if (response.ok) {
          const data = await response.json()
          setAlerts(data.alerts || [])
        }
      } catch (err) {
        console.error('獲取告警失敗:', err)
      }
    }

    // 初始載入
    fetchSystemStatus()
    fetchAlerts()

    // 每 30 秒更新一次
    const interval = setInterval(() => {
      fetchSystemStatus()
      fetchAlerts()
    }, 30000)

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

  const getAlertIcon = (level: string) => {
    switch (level) {
      case 'critical':
        return <AlertTriangle className="h-4 w-4 text-red-500" />
      case 'warning':
        return <AlertTriangle className="h-4 w-4 text-yellow-500" />
      case 'info':
        return <Eye className="h-4 w-4 text-blue-500" />
      default:
        return <Eye className="h-4 w-4 text-gray-500" />
    }
  }

  const handleNetworkControl = async (action: 'block' | 'unblock') => {
    try {
      const response = await fetch(`${apiBaseUrl}/api/v1/control/network`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ action }),
      })
      
      if (response.ok) {
        // 重新載入狀態
        window.location.reload()
      }
    } catch (err) {
      console.error('網路控制失敗:', err)
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-pandora-600 mx-auto mb-4"></div>
          <p className="text-gray-600">載入系統狀態中...</p>
        </div>
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
                <RefreshCw className="h-4 w-4 mr-2" />
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
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">
                Pandora Box Console IDS-IPS
              </h1>
              <p className="text-gray-600">
                智慧型入侵偵測與防護系統 - 即時監控儀表板
              </p>
            </div>
            <div className="flex items-center space-x-2">
              <Badge className="bg-green-100 text-green-800">
                <CheckCircle className="h-3 w-3 mr-1" />
                系統正常
              </Badge>
              <Button variant="outline" size="sm">
                <Settings className="h-4 w-4 mr-2" />
                設定
              </Button>
            </div>
          </div>
        </div>

        {/* 標籤導航 */}
        <div className="mb-8">
          <div className="border-b border-gray-200">
            <nav className="-mb-px flex space-x-8">
              {[
                { id: 'overview', label: '總覽', icon: BarChart3 },
                { id: 'security', label: '安全', icon: Shield },
                { id: 'network', label: '網路', icon: Network },
                { id: 'monitoring', label: '監控', icon: Activity }
              ].map(({ id, label, icon: Icon }) => (
                <button
                  key={id}
                  onClick={() => setActiveTab(id as any)}
                  className={`py-2 px-1 border-b-2 font-medium text-sm flex items-center ${
                    activeTab === id
                      ? 'border-pandora-500 text-pandora-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
                >
                  <Icon className="h-4 w-4 mr-2" />
                  {label}
                </button>
              ))}
            </nav>
          </div>
        </div>

        {/* 總覽標籤 */}
        {activeTab === 'overview' && (
          <>
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
                        版本: {systemStatus?.agent.version}
                      </p>
                      <p className="text-xs text-gray-600">
                        運行時間: {systemStatus?.agent.uptime}
                      </p>
                    </div>
                    {systemStatus && getStatusBadge(systemStatus.agent.status)}
                  </div>
                </CardContent>
              </Card>

              {/* 安全告警 */}
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">安全告警</CardTitle>
                  <Shield className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-red-600">
                    {systemStatus?.security.criticalAlerts || 0}
                  </div>
                  <p className="text-xs text-gray-600">
                    總告警: {systemStatus?.security.totalAlerts || 0}
                  </p>
                  <p className="text-xs text-gray-600">
                    已阻斷威脅: {systemStatus?.security.threatsBlocked || 0}
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
                    被阻斷 IP: {systemStatus?.network.blockedIPs || 0}
                  </p>
                  <p className="text-xs text-gray-600">
                    允許 IP: {systemStatus?.network.allowedIPs || 0}
                  </p>
                </CardContent>
              </Card>

              {/* 設備連接 */}
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">設備連接</CardTitle>
                  <Server className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-green-600">
                    {systemStatus?.devices.online || 0}
                  </div>
                  <p className="text-xs text-gray-600">
                    總設備: {systemStatus?.devices.total || 0}
                  </p>
                  <p className="text-xs text-gray-600">
                    離線: {systemStatus?.devices.offline || 0}
                  </p>
                </CardContent>
              </Card>
            </div>

            {/* 系統資源使用 */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
              {/* 系統指標 */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Cpu className="h-5 w-5 mr-2" />
                    系統資源使用
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-sm font-medium">CPU 使用率</span>
                      <span className="text-sm text-gray-600">
                        {systemStatus?.agent.cpuUsage || 0}%
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-pandora-600 h-2 rounded-full transition-all duration-300"
                        style={{ width: `${systemStatus?.agent.cpuUsage || 0}%` }}
                      ></div>
                    </div>
                    
                    <div className="flex justify-between items-center">
                      <span className="text-sm font-medium">記憶體使用</span>
                      <span className="text-sm text-gray-600">
                        {systemStatus?.agent.memoryUsage || 0}%
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-pandora-500 h-2 rounded-full transition-all duration-300"
                        style={{ width: `${systemStatus?.agent.memoryUsage || 0}%` }}
                      ></div>
                    </div>

                    <div className="flex justify-between items-center">
                      <span className="text-sm font-medium">磁碟使用</span>
                      <span className="text-sm text-gray-600">
                        {systemStatus?.agent.diskUsage || 0}%
                      </span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2">
                      <div 
                        className="bg-pandora-400 h-2 rounded-full transition-all duration-300"
                        style={{ width: `${systemStatus?.agent.diskUsage || 0}%` }}
                      ></div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* 網路流量 */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center">
                    <Network className="h-5 w-5 mr-2" />
                    網路流量統計
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex justify-between items-center">
                      <span className="text-sm font-medium">總流量</span>
                      <span className="text-sm text-gray-600">
                        {(systemStatus?.network.totalTraffic || 0).toLocaleString()} MB
                      </span>
                    </div>
                    
                    <div className="flex justify-between items-center">
                      <span className="text-sm font-medium">活躍連線</span>
                      <span className="text-sm text-gray-600">
                        {systemStatus?.console.activeConnections || 0}
                      </span>
                    </div>

                    <div className="flex justify-between items-center">
                      <span className="text-sm font-medium">錯誤率</span>
                      <span className="text-sm text-gray-600">
                        {systemStatus?.console.errorRate || 0}%
                      </span>
                    </div>

                    <div className="flex justify-between items-center">
                      <span className="text-sm font-medium">響應時間</span>
                      <span className="text-sm text-gray-600">
                        {systemStatus?.console.responseTime || 0}ms
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* 快速操作 */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Zap className="h-5 w-5 mr-2" />
                  快速操作
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                  <Button 
                    className="w-full" 
                    variant="outline"
                    onClick={() => window.location.reload()}
                  >
                    <RefreshCw className="h-4 w-4 mr-2" />
                    重新載入狀態
                  </Button>
                  
                  <Button 
                    className="w-full" 
                    variant="outline"
                    onClick={() => window.open('http://localhost:3000', '_blank')}
                  >
                    <BarChart3 className="h-4 w-4 mr-2" />
                    開啟 Grafana
                  </Button>
                  
                  <Button 
                    className="w-full" 
                    variant="outline"
                    onClick={() => window.open('http://localhost:9090', '_blank')}
                  >
                    <Activity className="h-4 w-4 mr-2" />
                    開啟 Prometheus
                  </Button>
                  
                  <Button 
                    className="w-full" 
                    variant="outline"
                    onClick={() => window.open('http://localhost:3100', '_blank')}
                  >
                    <Database className="h-4 w-4 mr-2" />
                    開啟 Loki
                  </Button>
                </div>
              </CardContent>
            </Card>
          </>
        )}

        {/* 安全標籤 */}
        {activeTab === 'security' && (
          <div className="space-y-6">
            {/* 安全統計 */}
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">嚴重告警</CardTitle>
                  <AlertTriangle className="h-4 w-4 text-red-500" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-red-600">
                    {systemStatus?.security.criticalAlerts || 0}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">警告</CardTitle>
                  <AlertTriangle className="h-4 w-4 text-yellow-500" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-yellow-600">
                    {systemStatus?.security.warningAlerts || 0}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">資訊</CardTitle>
                  <Eye className="h-4 w-4 text-blue-500" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-blue-600">
                    {systemStatus?.security.infoAlerts || 0}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">已阻斷威脅</CardTitle>
                  <Shield className="h-4 w-4 text-green-500" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-green-600">
                    {systemStatus?.security.threatsBlocked || 0}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* 最新告警 */}
            <Card>
              <CardHeader>
                <CardTitle>最新安全告警</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {alerts.slice(0, 5).map((alert) => (
                    <div key={alert.id} className="flex items-center justify-between p-3 border rounded-lg">
                      <div className="flex items-center space-x-3">
                        {getAlertIcon(alert.level)}
                        <div>
                          <p className="font-medium">{alert.message}</p>
                          <p className="text-sm text-gray-600">
                            {alert.source} • {new Date(alert.timestamp).toLocaleString()}
                          </p>
                        </div>
                      </div>
                      <Badge 
                        className={
                          alert.level === 'critical' ? 'bg-red-100 text-red-800' :
                          alert.level === 'warning' ? 'bg-yellow-100 text-yellow-800' :
                          'bg-blue-100 text-blue-800'
                        }
                      >
                        {alert.level}
                      </Badge>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {/* 網路標籤 */}
        {activeTab === 'network' && (
          <div className="space-y-6">
            {/* 網路控制 */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Network className="h-5 w-5 mr-2" />
                  網路控制
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <Button 
                    className="w-full" 
                    variant="destructive"
                    onClick={() => handleNetworkControl('block')}
                    disabled={systemStatus?.network.blocked}
                  >
                    <Lock className="h-4 w-4 mr-2" />
                    阻斷網路
                  </Button>
                  
                  <Button 
                    className="w-full" 
                    variant="default"
                    onClick={() => handleNetworkControl('unblock')}
                    disabled={!systemStatus?.network.blocked}
                  >
                    <Unlock className="h-4 w-4 mr-2" />
                    解除阻斷
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* 網路統計 */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="text-sm font-medium">總流量</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {(systemStatus?.network.totalTraffic || 0).toLocaleString()} MB
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-sm font-medium">被阻斷 IP</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-red-600">
                    {systemStatus?.network.blockedIPs || 0}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="text-sm font-medium">允許 IP</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold text-green-600">
                    {systemStatus?.network.allowedIPs || 0}
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        )}

        {/* 監控標籤 */}
        {activeTab === 'monitoring' && (
          <div className="space-y-6">
            {/* 監控服務狀態 */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Prometheus</CardTitle>
                  {systemStatus?.monitoring.prometheus ? (
                    <CheckCircle className="h-4 w-4 text-green-500" />
                  ) : (
                    <AlertTriangle className="h-4 w-4 text-red-500" />
                  )}
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {systemStatus?.monitoring.prometheus ? '正常' : '異常'}
                  </div>
                  <p className="text-xs text-gray-600">
                    指標收集服務
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Grafana</CardTitle>
                  {systemStatus?.monitoring.grafana ? (
                    <CheckCircle className="h-4 w-4 text-green-500" />
                  ) : (
                    <AlertTriangle className="h-4 w-4 text-red-500" />
                  )}
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {systemStatus?.monitoring.grafana ? '正常' : '異常'}
                  </div>
                  <p className="text-xs text-gray-600">
                    視覺化儀表板
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">Loki</CardTitle>
                  {systemStatus?.monitoring.loki ? (
                    <CheckCircle className="h-4 w-4 text-green-500" />
                  ) : (
                    <AlertTriangle className="h-4 w-4 text-red-500" />
                  )}
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {systemStatus?.monitoring.loki ? '正常' : '異常'}
                  </div>
                  <p className="text-xs text-gray-600">
                    日誌聚合服務
                  </p>
                </CardContent>
              </Card>

              <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                  <CardTitle className="text-sm font-medium">AlertManager</CardTitle>
                  {systemStatus?.monitoring.alertmanager ? (
                    <CheckCircle className="h-4 w-4 text-green-500" />
                  ) : (
                    <AlertTriangle className="h-4 w-4 text-red-500" />
                  )}
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">
                    {systemStatus?.monitoring.alertmanager ? '正常' : '異常'}
                  </div>
                  <p className="text-xs text-gray-600">
                    告警管理服務
                  </p>
                </CardContent>
              </Card>
            </div>

            {/* 監控連結 */}
            <Card>
              <CardHeader>
                <CardTitle>監控服務連結</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                  <Button 
                    className="w-full" 
                    variant="outline"
                    onClick={() => window.open('http://localhost:3000', '_blank')}
                  >
                    <BarChart3 className="h-4 w-4 mr-2" />
                    Grafana
                  </Button>
                  
                  <Button 
                    className="w-full" 
                    variant="outline"
                    onClick={() => window.open('http://localhost:9090', '_blank')}
                  >
                    <Activity className="h-4 w-4 mr-2" />
                    Prometheus
                  </Button>
                  
                  <Button 
                    className="w-full" 
                    variant="outline"
                    onClick={() => window.open('http://localhost:3100', '_blank')}
                  >
                    <Database className="h-4 w-4 mr-2" />
                    Loki
                  </Button>
                  
                  <Button 
                    className="w-full" 
                    variant="outline"
                    onClick={() => window.open('http://localhost:9093', '_blank')}
                  >
                    <AlertTriangle className="h-4 w-4 mr-2" />
                    AlertManager
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>
    </div>
  )
}