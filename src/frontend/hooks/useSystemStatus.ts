import { useState, useEffect } from 'react'
import { fetchSystemStatus } from '../services/api'

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

export function useSystemStatus(refreshInterval = 30000) {
  const [status, setStatus] = useState<SystemStatus | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let mounted = true
    let intervalId: NodeJS.Timeout

    const fetchData = async () => {
      try {
        setLoading(true)
        const data = await fetchSystemStatus()
        if (mounted) {
          setStatus(data)
          setError(null)
        }
      } catch (err) {
        if (mounted) {
          setError(err instanceof Error ? err.message : '獲取系統狀態失敗')
        }
      } finally {
        if (mounted) {
          setLoading(false)
        }
      }
    }

    fetchData()
    
    if (refreshInterval > 0) {
      intervalId = setInterval(fetchData, refreshInterval)
    }

    return () => {
      mounted = false
      if (intervalId) {
        clearInterval(intervalId)
      }
    }
  }, [refreshInterval])

  return { status, loading, error }
}

