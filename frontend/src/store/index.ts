import { create } from 'zustand'
import { persist } from 'zustand/middleware'

// 用户信息
export interface User {
  user_id: string
  username: string
  nickname: string
  avatar?: string
}

// 房间信息
export interface Room {
  room_id: string
  name: string
  description?: string
  creator_id: string
  user_count: number
  created_at: number
}

// 消息
export interface Message {
  msg_id: string
  room_id: string
  sender: User
  content: string
  msg_type: number
  timestamp: number
}

// 用户状态
interface UserState {
  user: User | null
  token: string | null
  isLoggedIn: boolean
  setUser: (user: User, token: string) => void
  logout: () => void
}

export const useUserStore = create<UserState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isLoggedIn: false,
      setUser: (user, token) => set({ user, token, isLoggedIn: true }),
      logout: () => set({ user: null, token: null, isLoggedIn: false }),
    }),
    {
      name: 'user-storage',
    }
  )
)

// 房间状态
interface RoomState {
  rooms: Room[]
  currentRoom: Room | null
  setRooms: (rooms: Room[]) => void
  setCurrentRoom: (room: Room | null) => void
  updateRoomUserCount: (roomId: string, count: number) => void
}

export const useRoomStore = create<RoomState>((set) => ({
  rooms: [],
  currentRoom: null,
  setRooms: (rooms) => set({ rooms }),
  setCurrentRoom: (room) => set({ currentRoom: room }),
  updateRoomUserCount: (roomId, count) =>
    set((state) => ({
      rooms: state.rooms.map((r) =>
        r.room_id === roomId ? { ...r, user_count: count } : r
      ),
    })),
}))

// 消息状态
interface MessageState {
  messages: Message[]
  hasMore: boolean
  addMessage: (message: Message) => void
  setMessages: (messages: Message[]) => void
  appendMessages: (messages: Message[]) => void
  clearMessages: () => void
}

export const useMessageStore = create<MessageState>((set) => ({
  messages: [],
  hasMore: false,
  addMessage: (message) =>
    set((state) => ({ messages: [...state.messages, message] })),
  setMessages: (messages) => set({ messages }),
  appendMessages: (messages) =>
    set((state) => ({ messages: [...messages, ...state.messages] })),
  clearMessages: () => set({ messages: [], hasMore: false }),
}))
