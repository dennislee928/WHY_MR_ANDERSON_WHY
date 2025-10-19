// Type definitions for the application

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
  details?: Record<string, any>
}

export interface NetworkControlRequest {
  action: 'block' | 'unblock'
}

export interface NetworkControlResponse {
  success: boolean
  message: string
  timestamp: number
}

export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: string
  message?: string
  timestamp?: number
}

export interface DashboardData {
  system_status: {
    status: 'healthy' | 'warning' | 'error'
    uptime: number
    device_connected: boolean
    network_blocked: boolean
    active_sessions: number
  }
  security_metrics: {
    threats_detected: number
    packets_inspected: number
    blocked_connections: number
    security_score: number
  }
  network_status: {
    inbound_traffic: number
    outbound_traffic: number
    latency: number
  }
}

export interface DeviceStatus {
  connected: boolean
  port: string
  type: string
  lastSeen: string
  firmware?: string
}

export interface ThreatInfo {
  id: string
  type: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  source: string
  target?: string
  timestamp: number
  blocked: boolean
  description: string
}

export interface WebSocketMessage {
  type: 'dashboard_update' | 'network_status' | 'device_status' | 'security_alert' | 'ping' | 'pong'
  data?: any
  timestamp?: number
}

