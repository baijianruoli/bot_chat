import React, { useEffect, useState } from 'react'
import {
  Card,
  List,
  Button,
  Modal,
  Form,
  Input,
  Tag,
  message,
  Space,
  Empty,
  Pagination,
} from 'antd'
import {
  PlusOutlined,
  TeamOutlined,
  EnterOutlined,
  LogoutOutlined,
} from '@ant-design/icons'
import { roomApi } from '../api'
import { useRoomStore, useUserStore } from '../store'
import { useNavigate } from 'react-router-dom'

const RoomList: React.FC = () => {
  const navigate = useNavigate()
  const { user } = useUserStore()
  const { rooms, setRooms, setCurrentRoom } = useRoomStore()
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)

  const fetchRooms = async (pageNum = 1) => {
    setLoading(true)
    try {
      const res: any = await roomApi.list({ page: pageNum, page_size: 10 })
      if (res.code === 0) {
        setRooms(res.data.rooms)
        setTotal(res.data.total)
      }
    } catch (error) {
      message.error('获取房间列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchRooms(page)
  }, [page])

  const handleCreateRoom = async (values: {
    name: string
    description?: string
  }) => {
    try {
      const res: any = await roomApi.create(values)
      if (res.code === 0) {
        message.success('创建成功')
        setModalVisible(false)
        form.resetFields()
        fetchRooms(page)
      } else {
        message.error(res.message || '创建失败')
      }
    } catch (error) {
      message.error('网络错误')
    }
  }

  const handleJoinRoom = async (roomId: string) => {
    try {
      const res: any = await roomApi.join({ room_id: roomId })
      if (res.code === 0) {
        message.success('加入成功')
        setCurrentRoom(res.data.room)
        navigate(`/chat/${roomId}`)
      } else {
        message.error(res.message || '加入失败')
      }
    } catch (error) {
      message.error('网络错误')
    }
  }

  return (
    <Card
      title="房间列表"
      extra={
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setModalVisible(true)}
        >
          创建房间
        </Button>
      }
    >
      <List
        loading={loading}
        dataSource={rooms}
        locale={{
          emptyText: <Empty description="暂无房间" />,
        }}
        renderItem={(room) => (
          <List.Item
            actions={[
              <Button
                type="primary"
                icon={<EnterOutlined />}
                onClick={() => handleJoinRoom(room.room_id)}
              >
                加入
              </Button>,
            ]}
          >
            <List.Item.Meta
              title={room.name}
              description={
                <Space direction="vertical" size="small">
                  <span>{room.description || '暂无描述'}</span>
                  <Space>
                    <Tag icon={<TeamOutlined />}>
                      {room.user_count} 人在线
                    </Tag>
                  </Space>
                </Space>
              }
            />
          </List.Item>
        )}
      />

      {total > 10 && (
        <Pagination
          current={page}
          total={total}
          pageSize={10}
          onChange={setPage}
          style={{ marginTop: 16, textAlign: 'right' }}
        />
      )}

      <Modal
        title="创建房间"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
      >
        <Form form={form} onFinish={handleCreateRoom} layout="vertical">
          <Form.Item
            name="name"
            label="房间名称"
            rules={[{ required: true, message: '请输入房间名称' }]}
          >
            <Input placeholder="房间名称" />
          </Form.Item>

          <Form.Item name="description" label="房间描述">
            <Input.TextArea placeholder="房间描述（可选）" rows={3} />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              创建
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

export default RoomList
