import React, { useEffect, useRef, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import {
  Card,
  Input,
  Button,
  List,
  Avatar,
  Space,
  Typography,
  Tag,
  Empty,
  message,
} from 'antd'
import {
  SendOutlined,
  UserOutlined,
  ArrowLeftOutlined,
  LoadingOutlined,
} from '@ant-design/icons'
import { useRoomStore, useMessageStore, useUserStore } from '../store'
import { messageApi } from '../api'
import dayjs from 'dayjs'

const { Text } = Typography

const Chat: React.FC = () => {
  const { roomId } = useParams<{ roomId: string }>()
  const navigate = useNavigate()
  const { user } = useUserStore()
  const { currentRoom, setCurrentRoom } = useRoomStore()
  const { messages, addMessage, setMessages, hasMore } = useMessageStore()
  const [inputValue, setInputValue] = useState('')
  const [loading, setLoading] = useState(false)
  const [sending, setSending] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  // 获取历史消息
  const fetchMessages = async () => {
    if (!roomId) return
    setLoading(true)
    try {
      const res: any = await messageApi.getHistory({
        room_id: roomId,
        limit: 50,
      })
      if (res.code === 0) {
        setMessages(res.data.messages)
      }
    } catch (error) {
      message.error('获取消息失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (roomId) {
      fetchMessages()
    }
  }, [roomId])

  // 滚动到底部
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = async () => {
    if (!inputValue.trim() || !roomId) return

    setSending(true)
    try {
      const res: any = await messageApi.send({
        room_id: roomId,
        content: inputValue.trim(),
        msg_type: 1,
      })

      if (res.code === 0) {
        addMessage(res.data.msg)
        setInputValue('')
      } else {
        message.error(res.message || '发送失败')
      }
    } catch (error) {
      message.error('网络错误')
    } finally {
      setSending(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const formatTime = (timestamp: number) => {
    return dayjs(timestamp).format('HH:mm')
  }

  return (
    <Card
      title={
        <Space>
          <Button
            icon={<ArrowLeftOutlined />}
            onClick={() => navigate('/rooms')}
          />
          <span>{currentRoom?.name || '聊天室'}</span>
          {currentRoom && (
            <Tag>{currentRoom.user_count} 人在线</Tag>
          )}
        </Space>
      }
      bodyStyle={{ padding: 0, height: 'calc(100vh - 180px)' }}
    >
      {/* 消息列表 */}
      <div
        style={{
          height: 'calc(100% - 60px)',
          overflowY: 'auto',
          padding: '16px',
        }}
      >
        {loading ? (
          <div style={{ textAlign: 'center', padding: '40px' }}>
            <LoadingOutlined style={{ fontSize: 24 }} />
          </div>
        ) : messages.length === 0 ? (
          <Empty description="暂无消息，开始聊天吧" />
        ) : (
          <List
            dataSource={messages}
            renderItem={(msg) => {
              const isMe = msg.sender.user_id === user?.user_id

              return (
                <List.Item
                  style={{
                    justifyContent: isMe ? 'flex-end' : 'flex-start',
                    padding: '8px 0',
                  }}
                >
                  <Space
                    align="start"
                    style={{
                      flexDirection: isMe ? 'row-reverse' : 'row',
                    }}
                  >
                    <Avatar
                      icon={<UserOutlined />}
                      src={msg.sender.avatar}
                    />
                    <div
                      style={{
                        maxWidth: '60%',
                        textAlign: isMe ? 'right' : 'left',
                      }}
                    >
                      <div style={{ marginBottom: 4 }}>
                        <Text strong>{msg.sender.nickname}</Text>
                        <Text type="secondary" style={{ marginLeft: 8, fontSize: 12 }}>
                          {formatTime(msg.timestamp)}
                        </Text>
                      </div>
                      <div
                        style={{
                          background: isMe ? '#1677ff' : '#f0f0f0',
                          color: isMe ? 'white' : 'inherit',
                          padding: '8px 12px',
                          borderRadius: 8,
                          display: 'inline-block',
                          wordBreak: 'break-word',
                        }}
                      >
                        {msg.content}
                      </div>
                    </div>
                  </Space>
                </List.Item>
              )
            }}
          />
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* 输入框 */}
      <div
        style={{
          padding: '12px 16px',
          borderTop: '1px solid #f0f0f0',
          display: 'flex',
          gap: 8,
        }}
      >
        <Input.TextArea
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="输入消息..."
          autoSize={{ minRows: 1, maxRows: 4 }}
          disabled={sending}
        />
        <Button
          type="primary"
          icon={<SendOutlined />}
          onClick={handleSend}
          loading={sending}
          disabled={!inputValue.trim()}
        >
          发送
        </Button>
      </div>
    </Card>
  )
}

export default Chat
