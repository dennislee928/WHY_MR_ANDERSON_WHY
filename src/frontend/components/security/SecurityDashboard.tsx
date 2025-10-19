import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card'
import { Button } from '../ui/button'
import { Badge } from '../ui/badge'
import { 
  Shield, 
  AlertTriangle, 
  Eye,
  Filter,
  Download,
  RefreshCw,
  CheckCircle,
  XCircle,
  Clock,
  TrendingUp,
  TrendingDown,
  Activity,
  Lock,
  Unlock,
  Target,
  Ban
} from 'lucide-react'

interface SecurityDashboardProps {
  apiBaseUrl: string
}

interface ThreatEvent {
  id: string
  timestamp: string
  type: string
  severity: 'critical' | 'high' | 'medium' | 'low'
  source_ip: string
  destination_ip: string
  port: number
  protocol: string
  action: 'blocked' | 'allowed' | 'logged'
  description: string
  rule_id: string
}

interface ThreatStats {
  total_threats: number
  blocked_threats: number
  active_threats: number
  resolved_threats: number
  threat_trend: number
  top_threat_types: Array<{type: string, count: number}>
  top_source_ips: Array<{ip: string, count: number}>
}

export default function SecurityDashboard({ apiBaseUrl }: SecurityDashboardProps) {
  const [threatEvents, setThreatEvents] = useState<ThreatEvent[]>([])
  const [threatStats, setThreatStats] = useState<ThreatStats | null>(null)
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<'all' | 'critical' | 'high' | 'medium' | 'low'>('all')
  const [timeRange, setTimeRange] = useState<'1h' | '24h' | '7d' | '30d'>('24h')

  useEffect(() => {
    fetchSecurityData()
    const interval = setInterval(fetchSecurityData, 30000)
    return () => clearInterval(interval)
  }, [apiBaseUrl, filter, timeRange])

  const fetchSecurityData = async () => {
    try {
      setLoading(true)
      
      // 獲取威脅事件
      const eventsResponse = await fetch(
        `${apiBaseUrl}/api/v1/security/threats?severity=${filter}&time_range=${timeRange}`
      )
      if (eventsResponse.ok) {
        const eventsData = await eventsResponse.json()
        setThreatEvents(eventsData.threats || [])
      }

      // 獲取威脅統計
      const statsResponse = await fetch(`${apiBaseUrl}/api/v1/security/stats`)
      if (statsResponse.ok) {
        const statsData = await statsResponse.json()
        setThreatStats(statsData)
      }
    } catch (err) {
      console.error('獲取安全數據失敗:', err)
    } finally {
      setLoading(false)
    }
  }

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'bg-red-100 text-red-800 border-red-300'
      case 'high':
        return 'bg-orange-100 text-orange-800 border-orange-300'
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300'
      case 'low':
        return 'bg-blue-100 text-blue-800 border-blue-300'
      default:
        return 'bg-gray-100 text-gray-800 border-gray-300'
    }
  }

  const getActionIcon = (action: string) => {
    switch (action) {
      case 'blocked':
        return <Ban className="h-4 w-4 text-red-500" />
      case 'allowed':
        return <CheckCircle className="h-4 w-4 text-green-500" />
      case 'logged':
        return <Eye className="h-4 w-4 text-blue-500" />
      default:
        return <Activity className="h-4 w-4 text-gray-500" />
    }
  }

  const exportReport = () => {
    // 實作報表匯出功能
    const csvContent = [
      ['時間戳', '威脅類型', '嚴重程度', '來源IP', '目標IP', '動作', '描述'].join(','),
      ...threatEvents.map(event => 
        [
          event.timestamp,
          event.type,
          event.severity,
          event.source_ip,
          event.destination_ip,
          event.action,
          event.description
        ].join(',')
      )
    ].join('\n')

    const blob = new Blob([csvContent], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `security_report_${new Date().toISOString()}.csv`
    a.click()
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      {/* 頁面標題 */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">安全監控</h1>
            <p className="text-gray-600">即時威脅偵測與安全事件管理</p>
          </div>
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="sm" onClick={fetchSecurityData}>
              <RefreshCw className="h-4 w-4 mr-2" />
              重新整理
            </Button>
            <Button variant="outline" size="sm" onClick={exportReport}>
              <Download className="h-4 w-4 mr-2" />
              匯出報表
            </Button>
          </div>
        </div>
      </div>

      {/* 威脅統計卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">總威脅數</CardTitle>
            <Shield className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{threatStats?.total_threats || 0}</div>
            <div className="flex items-center text-xs text-gray-600 mt-1">
              {(threatStats?.threat_trend || 0) > 0 ? (
                <>
                  <TrendingUp className="h-3 w-3 text-red-500 mr-1" />
                  <span className="text-red-500">+{threatStats?.threat_trend || 0}%</span>
                </>
              ) : (
                <>
                  <TrendingDown className="h-3 w-3 text-green-500 mr-1" />
                  <span className="text-green-500">{threatStats?.threat_trend || 0}%</span>
                </>
              )}
              <span className="ml-1">相較上週</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">已阻斷</CardTitle>
            <Lock className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {threatStats?.blocked_threats || 0}
            </div>
            <p className="text-xs text-gray-600 mt-1">自動阻斷惡意流量</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">活躍威脅</CardTitle>
            <AlertTriangle className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">
              {threatStats?.active_threats || 0}
            </div>
            <p className="text-xs text-gray-600 mt-1">需要立即處理</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">已解決</CardTitle>
            <CheckCircle className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {threatStats?.resolved_threats || 0}
            </div>
            <p className="text-xs text-gray-600 mt-1">威脅已消除</p>
          </CardContent>
        </Card>
      </div>

      {/* 過濾器和時間範圍 */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center space-x-2">
          <Filter className="h-5 w-5 text-gray-500" />
          <span className="text-sm font-medium text-gray-700">嚴重程度：</span>
          <div className="flex space-x-2">
            {['all', 'critical', 'high', 'medium', 'low'].map((severity) => (
              <Button
                key={severity}
                variant={filter === severity ? 'default' : 'outline'}
                size="sm"
                onClick={() => setFilter(severity as any)}
              >
                {severity === 'all' ? '全部' : severity}
              </Button>
            ))}
          </div>
        </div>

        <div className="flex items-center space-x-2">
          <Clock className="h-5 w-5 text-gray-500" />
          <span className="text-sm font-medium text-gray-700">時間範圍：</span>
          <div className="flex space-x-2">
            {[
              { value: '1h', label: '1小時' },
              { value: '24h', label: '24小時' },
              { value: '7d', label: '7天' },
              { value: '30d', label: '30天' }
            ].map(({ value, label }) => (
              <Button
                key={value}
                variant={timeRange === value ? 'default' : 'outline'}
                size="sm"
                onClick={() => setTimeRange(value as any)}
              >
                {label}
              </Button>
            ))}
          </div>
        </div>
      </div>

      {/* 威脅事件列表 */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Target className="h-5 w-5 mr-2" />
            威脅事件列表
          </CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-pandora-600 mx-auto"></div>
              <p className="text-gray-600 mt-2">載入中...</p>
            </div>
          ) : threatEvents.length === 0 ? (
            <div className="text-center py-8">
              <Shield className="h-12 w-12 text-gray-400 mx-auto mb-2" />
              <p className="text-gray-600">目前沒有威脅事件</p>
            </div>
          ) : (
            <div className="space-y-4">
              {threatEvents.map((event) => (
                <div
                  key={event.id}
                  className="p-4 border rounded-lg hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-start space-x-3 flex-1">
                      <div className="mt-1">
                        {getActionIcon(event.action)}
                      </div>
                      <div className="flex-1">
                        <div className="flex items-center space-x-2 mb-1">
                          <Badge className={`${getSeverityColor(event.severity)} border`}>
                            {event.severity}
                          </Badge>
                          <span className="text-sm font-medium text-gray-900">
                            {event.type}
                          </span>
                          <span className="text-xs text-gray-500">
                            {new Date(event.timestamp).toLocaleString()}
                          </span>
                        </div>
                        <p className="text-sm text-gray-700 mb-2">{event.description}</p>
                        <div className="flex items-center space-x-4 text-xs text-gray-600">
                          <span>來源: {event.source_ip}</span>
                          <span>目標: {event.destination_ip}</span>
                          <span>端口: {event.port}</span>
                          <span>協議: {event.protocol}</span>
                          <span>規則: {event.rule_id}</span>
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Badge
                        className={
                          event.action === 'blocked'
                            ? 'bg-red-100 text-red-800'
                            : event.action === 'allowed'
                            ? 'bg-green-100 text-green-800'
                            : 'bg-blue-100 text-blue-800'
                        }
                      >
                        {event.action}
                      </Badge>
                      <Button variant="outline" size="sm">
                        詳情
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* 威脅統計圖表 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-6">
        {/* 威脅類型分布 */}
        <Card>
          <CardHeader>
            <CardTitle>威脅類型分布</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {threatStats?.top_threat_types.map((item, index) => (
                <div key={index}>
                  <div className="flex justify-between text-sm mb-1">
                    <span className="font-medium">{item.type}</span>
                    <span className="text-gray-600">{item.count}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-pandora-600 h-2 rounded-full"
                      style={{
                        width: `${(item.count / (threatStats?.total_threats || 1)) * 100}%`
                      }}
                    ></div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* 高風險來源 IP */}
        <Card>
          <CardHeader>
            <CardTitle>高風險來源 IP</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {threatStats?.top_source_ips.map((item, index) => (
                <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                  <div>
                    <p className="font-medium text-sm">{item.ip}</p>
                    <p className="text-xs text-gray-600">{item.count} 次威脅事件</p>
                  </div>
                  <Button variant="destructive" size="sm">
                    <Ban className="h-4 w-4 mr-1" />
                    阻斷
                  </Button>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
