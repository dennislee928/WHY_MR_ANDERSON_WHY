import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card'
import { Button } from '../ui/button'
import { Badge } from '../ui/badge'
import { 
  Network, 
  Activity,
  TrendingUp,
  TrendingDown,
  Wifi,
  WifiOff,
  Lock,
  Unlock,
  Ban,
  CheckCircle,
  AlertCircle,
  Globe,
  Server,
  Zap,
  Clock
} from 'lucide-react'

interface NetworkDashboardProps {
  apiBaseUrl: string
}

interface NetworkStats {
  total_traffic: number
  inbound_traffic: number
  outbound_traffic: number
  packet_loss: number
  latency: number
  bandwidth_usage: number
  active_connections: number
  blocked_connections: number
}

interface BlockedIP {
  ip: string
  reason: string
  blocked_at: string
  expires_at: string
  threat_count: number
}

interface NetworkInterface {
  name: string
  status: 'up' | 'down'
  ip_address: string
  mac_address: string
  rx_bytes: number
  tx_bytes: number
  rx_packets: number
  tx_packets: number
}

export default function NetworkDashboard({ apiBaseUrl }: NetworkDashboardProps) {
  const [networkStats, setNetworkStats] = useState<NetworkStats | null>(null)
  const [blockedIPs, setBlockedIPs] = useState<BlockedIP[]>([])
  const [networkInterfaces, setNetworkInterfaces] = useState<NetworkInterface[]>([])
  const [loading, setLoading] = useState(true)
  const [networkBlocked, setNetworkBlocked] = useState(false)

  useEffect(() => {
    fetchNetworkData()
    const interval = setInterval(fetchNetworkData, 30000)
    return () => clearInterval(interval)
  }, [apiBaseUrl])

  const fetchNetworkData = async () => {
    try {
      setLoading(true)
      
      // 獲取網路統計
      const statsResponse = await fetch(`${apiBaseUrl}/api/v1/network/stats`)
      if (statsResponse.ok) {
        const data = await statsResponse.json()
        setNetworkStats(data)
      }

      // 獲取被阻斷的 IP
      const blockedResponse = await fetch(`${apiBaseUrl}/api/v1/network/blocked-ips`)
      if (blockedResponse.ok) {
        const data = await blockedResponse.json()
        setBlockedIPs(data.blocked_ips || [])
      }

      // 獲取網路介面
      const interfacesResponse = await fetch(`${apiBaseUrl}/api/v1/network/interfaces`)
      if (interfacesResponse.ok) {
        const data = await interfacesResponse.json()
        setNetworkInterfaces(data.interfaces || [])
      }

      // 獲取網路狀態
      const statusResponse = await fetch(`${apiBaseUrl}/api/v1/control/network/status`)
      if (statusResponse.ok) {
        const data = await statusResponse.json()
        setNetworkBlocked(data.blocked || false)
      }
    } catch (err) {
      console.error('獲取網路數據失敗:', err)
    } finally {
      setLoading(false)
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
        setNetworkBlocked(action === 'block')
        await fetchNetworkData()
      }
    } catch (err) {
      console.error('網路控制失敗:', err)
    }
  }

  const unblockIP = async (ip: string) => {
    try {
      const response = await fetch(`${apiBaseUrl}/api/v1/network/blocked-ips/${ip}`, {
        method: 'DELETE',
      })
      
      if (response.ok) {
        await fetchNetworkData()
      }
    } catch (err) {
      console.error('解除 IP 阻斷失敗:', err)
    }
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      {/* 頁面標題 */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">網路管理</h1>
            <p className="text-gray-600">網路流量監控與連線管理</p>
          </div>
          <div className="flex items-center space-x-2">
            <Badge className={networkBlocked ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'}>
              {networkBlocked ? (
                <>
                  <WifiOff className="h-3 w-3 mr-1" />
                  網路已阻斷
                </>
              ) : (
                <>
                  <Wifi className="h-3 w-3 mr-1" />
                  網路正常
                </>
              )}
            </Badge>
            <Button 
              variant={networkBlocked ? 'default' : 'destructive'}
              size="sm"
              onClick={() => handleNetworkControl(networkBlocked ? 'unblock' : 'block')}
            >
              {networkBlocked ? (
                <>
                  <Unlock className="h-4 w-4 mr-2" />
                  解除阻斷
                </>
              ) : (
                <>
                  <Lock className="h-4 w-4 mr-2" />
                  阻斷網路
                </>
              )}
            </Button>
          </div>
        </div>
      </div>

      {/* 網路統計卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">總流量</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {formatBytes(networkStats?.total_traffic || 0)}
            </div>
            <div className="flex items-center text-xs text-gray-600 mt-1">
              <TrendingUp className="h-3 w-3 text-green-500 mr-1" />
              <span>即時流量監控</span>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">連線數</CardTitle>
            <Globe className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {networkStats?.active_connections || 0}
            </div>
            <p className="text-xs text-gray-600 mt-1">活躍連線</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">延遲</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {networkStats?.latency || 0}ms
            </div>
            <p className="text-xs text-gray-600 mt-1">網路延遲</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">封包遺失</CardTitle>
            <AlertCircle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {networkStats?.packet_loss || 0}%
            </div>
            <p className="text-xs text-gray-600 mt-1">封包遺失率</p>
          </CardContent>
        </Card>
      </div>

      {/* 流量統計 */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle className="flex items-center">
            <Activity className="h-5 w-5 mr-2" />
            流量統計
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <div className="flex justify-between text-sm mb-2">
                <span className="font-medium">上傳流量</span>
                <span className="text-gray-600">
                  {formatBytes(networkStats?.outbound_traffic || 0)}
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-3">
                <div 
                  className="bg-blue-600 h-3 rounded-full transition-all"
                  style={{ 
                    width: `${((networkStats?.outbound_traffic || 0) / (networkStats?.total_traffic || 1)) * 100}%` 
                  }}
                ></div>
              </div>
            </div>

            <div>
              <div className="flex justify-between text-sm mb-2">
                <span className="font-medium">下載流量</span>
                <span className="text-gray-600">
                  {formatBytes(networkStats?.inbound_traffic || 0)}
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-3">
                <div 
                  className="bg-green-600 h-3 rounded-full transition-all"
                  style={{ 
                    width: `${((networkStats?.inbound_traffic || 0) / (networkStats?.total_traffic || 1)) * 100}%` 
                  }}
                ></div>
              </div>
            </div>
          </div>

          <div className="mt-6 pt-6 border-t">
            <div className="flex justify-between items-center">
              <div>
                <p className="text-sm font-medium text-gray-700">頻寬使用率</p>
                <p className="text-2xl font-bold text-pandora-600">
                  {networkStats?.bandwidth_usage || 0}%
                </p>
              </div>
              <div className="text-right">
                <p className="text-sm text-gray-600">已阻斷連線</p>
                <p className="text-2xl font-bold text-red-600">
                  {networkStats?.blocked_connections || 0}
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* 網路介面 */}
      <Card className="mb-8">
        <CardHeader>
          <CardTitle className="flex items-center">
            <Server className="h-5 w-5 mr-2" />
            網路介面
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {networkInterfaces.map((iface) => (
              <div
                key={iface.name}
                className="p-4 border rounded-lg hover:bg-gray-50 transition-colors"
              >
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center space-x-3">
                    {iface.status === 'up' ? (
                      <CheckCircle className="h-5 w-5 text-green-500" />
                    ) : (
                      <AlertCircle className="h-5 w-5 text-red-500" />
                    )}
                    <div>
                      <p className="font-medium text-gray-900">{iface.name}</p>
                      <p className="text-sm text-gray-600">{iface.ip_address}</p>
                    </div>
                  </div>
                  <Badge className={iface.status === 'up' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}>
                    {iface.status}
                  </Badge>
                </div>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                  <div>
                    <p className="text-gray-600">MAC 位址</p>
                    <p className="font-medium">{iface.mac_address}</p>
                  </div>
                  <div>
                    <p className="text-gray-600">接收</p>
                    <p className="font-medium">{formatBytes(iface.rx_bytes)}</p>
                  </div>
                  <div>
                    <p className="text-gray-600">發送</p>
                    <p className="font-medium">{formatBytes(iface.tx_bytes)}</p>
                  </div>
                  <div>
                    <p className="text-gray-600">封包</p>
                    <p className="font-medium">
                      ↓{iface.rx_packets.toLocaleString()} / ↑{iface.tx_packets.toLocaleString()}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* 被阻斷的 IP */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Ban className="h-5 w-5 mr-2" />
            被阻斷的 IP 位址
          </CardTitle>
        </CardHeader>
        <CardContent>
          {blockedIPs.length === 0 ? (
            <div className="text-center py-8">
              <CheckCircle className="h-12 w-12 text-gray-400 mx-auto mb-2" />
              <p className="text-gray-600">目前沒有被阻斷的 IP</p>
            </div>
          ) : (
            <div className="space-y-3">
              {blockedIPs.map((item) => (
                <div
                  key={item.ip}
                  className="p-4 border rounded-lg hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex-1">
                      <div className="flex items-center space-x-3 mb-2">
                        <Ban className="h-4 w-4 text-red-500" />
                        <span className="font-medium text-gray-900">{item.ip}</span>
                        <Badge className="bg-red-100 text-red-800">
                          {item.threat_count} 次威脅
                        </Badge>
                      </div>
                      <p className="text-sm text-gray-700 mb-2">{item.reason}</p>
                      <div className="flex items-center space-x-4 text-xs text-gray-600">
                        <span className="flex items-center">
                          <Clock className="h-3 w-3 mr-1" />
                          阻斷時間: {new Date(item.blocked_at).toLocaleString()}
                        </span>
                        <span className="flex items-center">
                          <Clock className="h-3 w-3 mr-1" />
                          到期時間: {new Date(item.expires_at).toLocaleString()}
                        </span>
                      </div>
                    </div>
                    <Button 
                      variant="outline" 
                      size="sm"
                      onClick={() => unblockIP(item.ip)}
                    >
                      <Unlock className="h-4 w-4 mr-1" />
                      解除阻斷
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
