import React, { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card'
import { Button } from '../ui/button'
import { Badge } from '../ui/badge'
import { 
  Settings, 
  Bell,
  Shield,
  Network,
  Database,
  Users,
  Key,
  Mail,
  Smartphone,
  Clock,
  Globe,
  Save,
  RefreshCw
} from 'lucide-react'

interface SettingsDashboardProps {
  apiBaseUrl: string
}

export default function SettingsDashboard({ apiBaseUrl }: SettingsDashboardProps) {
  const [activeTab, setActiveTab] = useState<'general' | 'security' | 'notifications' | 'network'>('general')
  const [settings, setSettings] = useState({
    general: {
      system_name: 'Pandora Box Console IDS-IPS',
      language: 'zh-TW',
      timezone: 'Asia/Taipei',
      auto_update: true,
      log_level: 'info',
    },
    security: {
      two_factor_auth: false,
      session_timeout: 30,
      password_policy: 'strong',
      ip_whitelist_enabled: false,
      max_login_attempts: 5,
    },
    notifications: {
      email_enabled: true,
      email_address: 'admin@pandora-ids.com',
      slack_enabled: false,
      slack_webhook: '',
      alert_threshold: 'medium',
      digest_frequency: 'daily',
    },
    network: {
      auto_block_enabled: true,
      block_duration: 24,
      max_connections: 1000,
      rate_limiting: true,
      ddos_protection: true,
    },
  })

  const handleSave = async () => {
    // 實作保存設定
    console.log('保存設定:', settings)
  }

  const handleReset = () => {
    // 實作重置設定
    console.log('重置設定')
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      {/* 頁面標題 */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900 mb-2">系統設定</h1>
            <p className="text-gray-600">配置系統參數和偏好設定</p>
          </div>
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="sm" onClick={handleReset}>
              <RefreshCw className="h-4 w-4 mr-2" />
              重置
            </Button>
            <Button size="sm" onClick={handleSave}>
              <Save className="h-4 w-4 mr-2" />
              儲存變更
            </Button>
          </div>
        </div>
      </div>

      {/* 標籤導航 */}
      <div className="mb-8">
        <div className="border-b border-gray-200">
          <nav className="-mb-px flex space-x-8">
            {[
              { id: 'general', label: '一般設定', icon: Settings },
              { id: 'security', label: '安全性', icon: Shield },
              { id: 'notifications', label: '通知', icon: Bell },
              { id: 'network', label: '網路', icon: Network }
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

      {/* 一般設定 */}
      {activeTab === 'general' && (
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>系統資訊</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  系統名稱
                </label>
                <input
                  type="text"
                  value={settings.general.system_name}
                  onChange={(e) => setSettings({
                    ...settings,
                    general: { ...settings.general, system_name: e.target.value }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  <Globe className="h-4 w-4 inline mr-1" />
                  語言
                </label>
                <select
                  value={settings.general.language}
                  onChange={(e) => setSettings({
                    ...settings,
                    general: { ...settings.general, language: e.target.value }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                >
                  <option value="zh-TW">繁體中文</option>
                  <option value="en-US">English</option>
                  <option value="ja-JP">日本語</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  <Clock className="h-4 w-4 inline mr-1" />
                  時區
                </label>
                <select
                  value={settings.general.timezone}
                  onChange={(e) => setSettings({
                    ...settings,
                    general: { ...settings.general, timezone: e.target.value }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                >
                  <option value="Asia/Taipei">台北 (GMT+8)</option>
                  <option value="UTC">UTC (GMT+0)</option>
                  <option value="America/New_York">紐約 (GMT-5)</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  日誌等級
                </label>
                <select
                  value={settings.general.log_level}
                  onChange={(e) => setSettings({
                    ...settings,
                    general: { ...settings.general, log_level: e.target.value }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                >
                  <option value="debug">Debug</option>
                  <option value="info">Info</option>
                  <option value="warning">Warning</option>
                  <option value="error">Error</option>
                </select>
              </div>

              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="auto_update"
                  checked={settings.general.auto_update}
                  onChange={(e) => setSettings({
                    ...settings,
                    general: { ...settings.general, auto_update: e.target.checked }
                  })}
                  className="h-4 w-4 text-pandora-600 focus:ring-pandora-500 border-gray-300 rounded"
                />
                <label htmlFor="auto_update" className="ml-2 block text-sm text-gray-700">
                  啟用自動更新
                </label>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* 安全性設定 */}
      {activeTab === 'security' && (
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Key className="h-5 w-5 mr-2" />
                認證設定
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-gray-900">雙因子認證</p>
                  <p className="text-sm text-gray-600">使用 OTP 應用程式進行額外驗證</p>
                </div>
                <input
                  type="checkbox"
                  checked={settings.security.two_factor_auth}
                  onChange={(e) => setSettings({
                    ...settings,
                    security: { ...settings.security, two_factor_auth: e.target.checked }
                  })}
                  className="h-4 w-4 text-pandora-600 focus:ring-pandora-500 border-gray-300 rounded"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  會話逾時（分鐘）
                </label>
                <input
                  type="number"
                  value={settings.security.session_timeout}
                  onChange={(e) => setSettings({
                    ...settings,
                    security: { ...settings.security, session_timeout: parseInt(e.target.value) }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  密碼政策
                </label>
                <select
                  value={settings.security.password_policy}
                  onChange={(e) => setSettings({
                    ...settings,
                    security: { ...settings.security, password_policy: e.target.value }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                >
                  <option value="weak">弱（6 字符）</option>
                  <option value="medium">中（8 字符 + 數字）</option>
                  <option value="strong">強（12 字符 + 符號）</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  最大登入嘗試次數
                </label>
                <input
                  type="number"
                  value={settings.security.max_login_attempts}
                  onChange={(e) => setSettings({
                    ...settings,
                    security: { ...settings.security, max_login_attempts: parseInt(e.target.value) }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                />
              </div>

              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="ip_whitelist"
                  checked={settings.security.ip_whitelist_enabled}
                  onChange={(e) => setSettings({
                    ...settings,
                    security: { ...settings.security, ip_whitelist_enabled: e.target.checked }
                  })}
                  className="h-4 w-4 text-pandora-600 focus:ring-pandora-500 border-gray-300 rounded"
                />
                <label htmlFor="ip_whitelist" className="ml-2 block text-sm text-gray-700">
                  啟用 IP 白名單
                </label>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* 通知設定 */}
      {activeTab === 'notifications' && (
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Mail className="h-5 w-5 mr-2" />
                電子郵件通知
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="email_enabled"
                  checked={settings.notifications.email_enabled}
                  onChange={(e) => setSettings({
                    ...settings,
                    notifications: { ...settings.notifications, email_enabled: e.target.checked }
                  })}
                  className="h-4 w-4 text-pandora-600 focus:ring-pandora-500 border-gray-300 rounded"
                />
                <label htmlFor="email_enabled" className="ml-2 block text-sm text-gray-700">
                  啟用電子郵件通知
                </label>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  電子郵件地址
                </label>
                <input
                  type="email"
                  value={settings.notifications.email_address}
                  onChange={(e) => setSettings({
                    ...settings,
                    notifications: { ...settings.notifications, email_address: e.target.value }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  告警門檻
                </label>
                <select
                  value={settings.notifications.alert_threshold}
                  onChange={(e) => setSettings({
                    ...settings,
                    notifications: { ...settings.notifications, alert_threshold: e.target.value }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                >
                  <option value="low">低（所有告警）</option>
                  <option value="medium">中（中等以上）</option>
                  <option value="high">高（僅嚴重告警）</option>
                  <option value="critical">嚴重（僅關鍵告警）</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  摘要頻率
                </label>
                <select
                  value={settings.notifications.digest_frequency}
                  onChange={(e) => setSettings({
                    ...settings,
                    notifications: { ...settings.notifications, digest_frequency: e.target.value }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                >
                  <option value="realtime">即時</option>
                  <option value="hourly">每小時</option>
                  <option value="daily">每日</option>
                  <option value="weekly">每週</option>
                </select>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Smartphone className="h-5 w-5 mr-2" />
                Slack 通知
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="slack_enabled"
                  checked={settings.notifications.slack_enabled}
                  onChange={(e) => setSettings({
                    ...settings,
                    notifications: { ...settings.notifications, slack_enabled: e.target.checked }
                  })}
                  className="h-4 w-4 text-pandora-600 focus:ring-pandora-500 border-gray-300 rounded"
                />
                <label htmlFor="slack_enabled" className="ml-2 block text-sm text-gray-700">
                  啟用 Slack 通知
                </label>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Webhook URL
                </label>
                <input
                  type="text"
                  value={settings.notifications.slack_webhook}
                  onChange={(e) => setSettings({
                    ...settings,
                    notifications: { ...settings.notifications, slack_webhook: e.target.value }
                  })}
                  placeholder="https://hooks.slack.com/services/..."
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                />
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* 網路設定 */}
      {activeTab === 'network' && (
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>網路保護</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-gray-900">自動阻斷</p>
                  <p className="text-sm text-gray-600">自動阻斷偵測到的威脅 IP</p>
                </div>
                <input
                  type="checkbox"
                  checked={settings.network.auto_block_enabled}
                  onChange={(e) => setSettings({
                    ...settings,
                    network: { ...settings.network, auto_block_enabled: e.target.checked }
                  })}
                  className="h-4 w-4 text-pandora-600 focus:ring-pandora-500 border-gray-300 rounded"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  阻斷持續時間（小時）
                </label>
                <input
                  type="number"
                  value={settings.network.block_duration}
                  onChange={(e) => setSettings({
                    ...settings,
                    network: { ...settings.network, block_duration: parseInt(e.target.value) }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  最大連線數
                </label>
                <input
                  type="number"
                  value={settings.network.max_connections}
                  onChange={(e) => setSettings({
                    ...settings,
                    network: { ...settings.network, max_connections: parseInt(e.target.value) }
                  })}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-pandora-500"
                />
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-gray-900">流量限制</p>
                  <p className="text-sm text-gray-600">限制單一 IP 的請求速率</p>
                </div>
                <input
                  type="checkbox"
                  checked={settings.network.rate_limiting}
                  onChange={(e) => setSettings({
                    ...settings,
                    network: { ...settings.network, rate_limiting: e.target.checked }
                  })}
                  className="h-4 w-4 text-pandora-600 focus:ring-pandora-500 border-gray-300 rounded"
                />
              </div>

              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium text-gray-900">DDoS 保護</p>
                  <p className="text-sm text-gray-600">啟用分散式阻斷服務攻擊防護</p>
                </div>
                <input
                  type="checkbox"
                  checked={settings.network.ddos_protection}
                  onChange={(e) => setSettings({
                    ...settings,
                    network: { ...settings.network, ddos_protection: e.target.checked }
                  })}
                  className="h-4 w-4 text-pandora-600 focus:ring-pandora-500 border-gray-300 rounded"
                />
              </div>
            </CardContent>
          </Card>
        </div>
      )}
    </div>
  )
}
