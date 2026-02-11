import axios from 'axios'
import { User, Room, Message } from '../store'

// API 基础配置
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 添加 token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器 - 统一错误处理
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 通用响应结构
interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

// ==================== 用户相关 API ====================

export interface RegisterReq {
  username: string
  password: string
  nickname: string
}

export interface RegisterResp {
  user_id: string
}

export interface LoginReq {
  username: string
  password: string
}

export interface LoginResp {
  token: string
  user: User
}

export const userApi = {
  register: (data: RegisterReq) =>
    api.post<ApiResponse<RegisterResp>>('/register', data),
  
  login: (data: LoginReq) =>
    api.post<ApiResponse<LoginResp>>('/login', data),
}

// ==================== 房间相关 API ====================

export interface CreateRoomReq {
  name: string
  description?: string
}

export interface CreateRoomResp {
  room: Room
}

export interface ListRoomsReq {
  page?: number
  page_size?: number
}

export interface ListRoomsResp {
  rooms: Room[]
  total: number
}

export interface JoinRoomReq {
  room_id: string
}

export interface JoinRoomResp {
  room: Room
}

export interface LeaveRoomReq {
  room_id: string
}

export const roomApi = {
  create: (data: CreateRoomReq) =>
    api.post<ApiResponse<CreateRoomResp>>('/rooms', data),
  
  list: (params?: ListRoomsReq) =>
    api.get<ApiResponse<ListRoomsResp>>('/rooms', { params }),
  
  join: (data: JoinRoomReq) =>
    api.post<ApiResponse<JoinRoomResp>>(`/rooms/${data.room_id}/join`, {}),
  
  leave: (data: LeaveRoomReq) =>
    api.post<ApiResponse<void>>(`/rooms/${data.room_id}/leave`, {}),
}

// ==================== 消息相关 API ====================

export interface SendMessageReq {
  room_id: string
  content: string
  msg_type?: number
}

export interface SendMessageResp {
  msg: Message
}

export interface GetHistoryReq {
  room_id: string
  before_time?: number
  limit?: number
}

export interface GetHistoryResp {
  messages: Message[]
  has_more: boolean
}

export const messageApi = {
  send: (data: SendMessageReq) =>
    api.post<ApiResponse<SendMessageResp>>('/messages', data),
  
  getHistory: (params: GetHistoryReq) =>
    api.get<ApiResponse<GetHistoryResp>>('/messages', { params }),
}

export default api
