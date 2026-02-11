import { useEffect, useRef, useCallback } from 'react'
import { useUserStore, useMessageStore, Message } from '../store'

const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8888/ws'

export const useWebSocket = (roomId: string | undefined) => {
  const { user, token } = useUserStore()
  const { addMessage } = useMessageStore()
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>()

  const connect = useCallback(() => {
    if (!roomId || !user) return

    const wsUrl = `${WS_URL}?user_id=${user.user_id}&token=${token}`
    const ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      console.log('WebSocket connected')
      // 发送加入房间消息
      ws.send(JSON.stringify({
        type: 'join',
        room_id: roomId,
      }))
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        
        switch (data.type) {
          case 'message':
            // 新消息
            addMessage(data.data as Message)
            break
          case 'join':
            // 用户加入
            console.log('User joined:', data.data)
            break
          case 'leave':
            // 用户离开
            console.log('User left:', data.data)
            break
          case 'online_count':
            // 在线人数更新
            console.log('Online count:', data.data.count)
            break
          default:
            console.log('Unknown message type:', data.type)
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    ws.onclose = () => {
      console.log('WebSocket disconnected')
      // 尝试重连
      reconnectTimeoutRef.current = setTimeout(() => {
        connect()
      }, 3000)
    }

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
    }

    wsRef.current = ws
  }, [roomId, user, token, addMessage])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
  }, [])

  const sendMessage = useCallback((content: string, msgType: number = 1) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      console.error('WebSocket not connected')
      return false
    }

    wsRef.current.send(JSON.stringify({
      type: 'message',
      room_id: roomId,
      content,
      msg_type: msgType,
    }))

    return true
  }, [roomId])

  useEffect(() => {
    connect()
    return () => disconnect()
  }, [connect, disconnect])

  return {
    sendMessage,
    isConnected: wsRef.current?.readyState === WebSocket.OPEN,
  }
}
