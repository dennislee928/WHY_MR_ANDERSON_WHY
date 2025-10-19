// API Service Layer
// 統一的API請求處理

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public data?: any
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

async function fetchAPI<T = any>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`
  
  const config: RequestInit = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  }

  try {
    const response = await fetch(url, config)
    
    if (!response.ok) {
      const error = await response.json().catch(() => ({}))
      throw new ApiError(
        error.message || `HTTP ${response.status}: ${response.statusText}`,
        response.status,
        error
      )
    }

    return await response.json()
  } catch (error) {
    if (error instanceof ApiError) {
      throw error
    }
    throw new ApiError('Network error or server unreachable', 0)
  }
}

// System Status APIs
export async function fetchSystemStatus() {
  return fetchAPI('/api/v1/status')
}

export async function fetchDashboardData() {
  return fetchAPI('/api/v1/dashboard')
}

// Network Control APIs
export async function blockNetwork() {
  return fetchAPI('/api/v1/control/network', {
    method: 'POST',
    body: JSON.stringify({ action: 'block' }),
  })
}

export async function unblockNetwork() {
  return fetchAPI('/api/v1/control/network', {
    method: 'POST',
    body: JSON.stringify({ action: 'unblock' }),
  })
}

// Events APIs
export async function fetchEvents(limit = 10, type?: string) {
  const params = new URLSearchParams({
    limit: limit.toString(),
    ...(type && { type }),
  })
  return fetchAPI(`/api/v1/events?${params}`)
}

export async function fetchEventById(id: string) {
  return fetchAPI(`/api/v1/events/${id}`)
}

// Security APIs
export async function fetchSecurityMetrics() {
  return fetchAPI('/api/v1/security/metrics')
}

export async function fetchThreats() {
  return fetchAPI('/api/v1/security/threats')
}

// Device APIs
export async function fetchDeviceStatus() {
  return fetchAPI('/api/v1/device/status')
}

// WebSocket Helper
export function createWebSocket(clientId?: string): WebSocket {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = process.env.NEXT_PUBLIC_WS_URL || `${protocol}//${window.location.host}/ws`
  const id = clientId || `client_${Date.now()}`
  
  return new WebSocket(`${wsUrl}?client_id=${id}`)
}

export default {
  fetchSystemStatus,
  fetchDashboardData,
  blockNetwork,
  unblockNetwork,
  fetchEvents,
  fetchEventById,
  fetchSecurityMetrics,
  fetchThreats,
  fetchDeviceStatus,
  createWebSocket,
}

