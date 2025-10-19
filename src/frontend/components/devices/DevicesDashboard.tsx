import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card'
import { Button } from '../ui/button'
import { Badge } from '../ui/badge'
import { 
  Server, 
  Activity,
  CheckCircle,
  XCircle,
  RefreshCw,
  Settings,
  Power,
  Info,
  AlertCircle,
  Usb,
  Network as NetworkIcon,
  Wifi,
  HardDrive
} from 'lucide-react'

interface DevicesDashboardProps {
  apiBaseUrl: string
}

interface Device {
  id: string
  name: string
  type: string
  status: 'online' | 'offline' | 'error'
  port: string
  baud_rate?: number
  ip_address?: string
  mac_address?: string
  last_seen: string
  uptime: string
  rx_bytes: number
  tx_bytes: number
  error_count: number
}

export default function DevicesDashboard({ apiBaseUrl }: DevicesDashboardProps) {
  const [devices, setDevices] = useState<Device[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedDevice, setSelectedDevice] = useState<string | null>(null)

  useEffect(() => {
    fetchDevices()
    const interval = setInterval(fetchDevices, 30000)
    return () => clearInterval(interval)
  }, [apiBaseUrl])

  const fetchDevices = async () => {
    try {
      setLoading(true)
      const response = await fetch(`${apiBaseUrl}/api/v1/devices`)
      if (response.ok) {
        const data = await response.json()
        setDevices(data.devices || [])
      }
    } catch (err) {
      console.error('獲取設備列表失敗:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleRestartDevice = async (deviceId: string) => {
    try {
      const response = await fetch(`${apiBaseUrl}/api/v1/devices/${deviceId}/restart`, {
        method: 'POST',
      })
      if (response.ok) {
        await fetchDevices()
      }
    } catch (err) {
      console.error('重啟設備失敗:', err)
    }
  }

  const getDeviceIcon = (type: string) => {
    switch (type) {
      case 'serial':
        return <Usb className="h-6 w-6 text-pandora-600" />
      case 'network':
        return <NetworkIcon className="h-6 w-6 text-pandora-600" />
      case 'sensor':
        return <Activity className="h-6 w-6 text-pandora-600" />
      default:
        return <Server className="h-6 w-6 text-pandora-600" />
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online':
        return <CheckCircle className="h-5 w-5 text-green-500" />
      case 'offline':
        return <XCircle className="h-5 w-5 text-red-500" />
      case 'error':
        return <AlertCircle className="h-5 w-5 text-yellow-500" />
      default:
        return <Info className="h-5 w-5 text-gray-500" />
    }
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
  }

  const onlineDevices = devices.filter(d => d.status === 'online').length
  const offlineDevices = devices.filter(d => d.status === 'offline').length

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      {/* 頁面標題 */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">設備管理</h1>
            <p className="text-gray-600">監控和管理所有連接的設備</p>
          </div>
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="sm" onClick={fetchDevices}>
              <RefreshCw className="h-4 w-4 mr-2" />
              重新整理
            </Button>
          </div>
        </div>
      </div>

      {/* 統計卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">總設備數</CardTitle>
            <Server className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{devices.length}</div>
            <p className="text-xs text-gray-600 mt-1">已註冊設備</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">線上設備</CardTitle>
            <CheckCircle className="h-4 w-4 text-green-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">{onlineDevices}</div>
            <p className="text-xs text-gray-600 mt-1">正常運行中</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">離線設備</CardTitle>
            <XCircle className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{offlineDevices}</div>
            <p className="text-xs text-gray-600 mt-1">需要檢查</p>
          </CardContent>
        </Card>
      </div>

      {/* 設備列表 */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Server className="h-5 w-5 mr-2" />
            設備列表
          </CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-pandora-600 mx-auto"></div>
              <p className="text-gray-600 mt-2">載入中...</p>
            </div>
          ) : devices.length === 0 ? (
            <div className="text-center py-8">
              <Server className="h-12 w-12 text-gray-400 mx-auto mb-2" />
              <p className="text-gray-600">目前沒有設備</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {devices.map((device) => (
                <Card
                  key={device.id}
                  className={`hover:shadow-lg transition-shadow ${
                    selectedDevice === device.id ? 'border-pandora-500 border-2' : ''
                  }`}
                >
                  <CardHeader>
                    <div className="flex items-start justify-between">
                      <div className="flex items-center space-x-3">
                        {getDeviceIcon(device.type)}
                        <div>
                          <h3 className="font-semibold text-gray-900">{device.name}</h3>
                          <p className="text-sm text-gray-600">{device.type}</p>
                        </div>
                      </div>
                      {getStatusIcon(device.status)}
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3">
                      {/* 狀態 Badge */}
                      <div className="flex items-center space-x-2">
                        <Badge
                          className={
                            device.status === 'online'
                              ? 'bg-green-100 text-green-800'
                              : device.status === 'offline'
                              ? 'bg-red-100 text-red-800'
                              : 'bg-yellow-100 text-yellow-800'
                          }
                        >
                          {device.status}
                        </Badge>
                        {device.error_count > 0 && (
                          <Badge className="bg-orange-100 text-orange-800">
                            {device.error_count} 錯誤
                          </Badge>
                        )}
                      </div>

                      {/* 設備資訊 */}
                      <div className="space-y-2 text-sm">
                        <div className="flex justify-between">
                          <span className="text-gray-600">端口:</span>
                          <span className="font-medium">{device.port}</span>
                        </div>
                        {device.baud_rate && (
                          <div className="flex justify-between">
                            <span className="text-gray-600">鮑率:</span>
                            <span className="font-medium">{device.baud_rate}</span>
                          </div>
                        )}
                        {device.ip_address && (
                          <div className="flex justify-between">
                            <span className="text-gray-600">IP:</span>
                            <span className="font-medium">{device.ip_address}</span>
                          </div>
                        )}
                        {device.mac_address && (
                          <div className="flex justify-between">
                            <span className="text-gray-600">MAC:</span>
                            <span className="font-medium text-xs">{device.mac_address}</span>
                          </div>
                        )}
                        <div className="flex justify-between">
                          <span className="text-gray-600">運行時間:</span>
                          <span className="font-medium">{device.uptime}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-gray-600">接收:</span>
                          <span className="font-medium">{formatBytes(device.rx_bytes)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-gray-600">發送:</span>
                          <span className="font-medium">{formatBytes(device.tx_bytes)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-gray-600">最後活動:</span>
                          <span className="font-medium text-xs">
                            {new Date(device.last_seen).toLocaleString()}
                          </span>
                        </div>
                      </div>

                      {/* 操作按鈕 */}
                      <div className="flex space-x-2 pt-3 border-t">
                        <Button
                          variant="outline"
                          size="sm"
                          className="flex-1"
                          onClick={() => setSelectedDevice(device.id)}
                        >
                          <Info className="h-4 w-4 mr-1" />
                          詳情
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleRestartDevice(device.id)}
                          disabled={device.status === 'offline'}
                        >
                          <Power className="h-4 w-4" />
                        </Button>
                        <Button variant="outline" size="sm">
                          <Settings className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
