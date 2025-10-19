import { useState, useEffect, useCallback, useRef } from 'react'
import { createWebSocket } from '../services/api'

interface WebSocketMessage {
  type: string
  data?: any
}

interface UseWebSocketOptions {
  onMessage?: (message: WebSocketMessage) => void
  onConnect?: () => void
  onDisconnect?: () => void
  onError?: (error: Event) => void
  reconnectInterval?: number
  heartbeatInterval?: number
}

export function useWebSocket(options: UseWebSocketOptions = {}) {
  const {
    onMessage,
    onConnect,
    onDisconnect,
    onError,
    reconnectInterval = 5000,
    heartbeatInterval = 30000,
  } = options

  const [connected, setConnected] = useState(false)
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimerRef = useRef<NodeJS.Timeout>()
  const heartbeatTimerRef = useRef<NodeJS.Timeout>()

  const connect = useCallback(() => {
    try {
      const ws = createWebSocket()
      wsRef.current = ws

      ws.onopen = () => {
        setConnected(true)
        onConnect?.()
        
        // 開始心跳
        heartbeatTimerRef.current = setInterval(() => {
          if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type: 'ping' }))
          }
        }, heartbeatInterval)
      }

      ws.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data) as WebSocketMessage
          setLastMessage(message)
          onMessage?.(message)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
        onError?.(error)
      }

      ws.onclose = () => {
        setConnected(false)
        onDisconnect?.()
        
        // 清理心跳
        if (heartbeatTimerRef.current) {
          clearInterval(heartbeatTimerRef.current)
        }
        
        // 自動重連
        reconnectTimerRef.current = setTimeout(() => {
          console.log('Attempting to reconnect WebSocket...')
          connect()
        }, reconnectInterval)
      }
    } catch (error) {
      console.error('Failed to create WebSocket:', error)
    }
  }, [onConnect, onDisconnect, onError, onMessage, reconnectInterval, heartbeatInterval])

  const disconnect = useCallback(() => {
    if (reconnectTimerRef.current) {
      clearTimeout(reconnectTimerRef.current)
    }
    if (heartbeatTimerRef.current) {
      clearInterval(heartbeatTimerRef.current)
    }
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
  }, [])

  const sendMessage = useCallback((message: WebSocketMessage) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message))
    } else {
      console.warn('WebSocket is not connected')
    }
  }, [])

  useEffect(() => {
    connect()
    return () => {
      disconnect()
    }
  }, [connect, disconnect])

  return {
    connected,
    lastMessage,
    sendMessage,
    reconnect: connect,
    disconnect,
  }
}

