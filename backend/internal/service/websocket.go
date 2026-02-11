package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocket 连接管理器
type WSManager struct {
	clients    map[string]*WSClient // userID -> client
	rooms      map[string]map[string]*WSClient // roomID -> userID -> client
	broadcast  chan *WSMessage
	register   chan *WSClient
	unregister chan *WSClient
	mu         sync.RWMutex
}

// WSClient WebSocket 客户端
type WSClient struct {
	manager *WSManager
	conn    *websocket.Conn
	send    chan []byte
	userID  string
	roomID  string
}

// WSMessage WebSocket 消息
type WSMessage struct {
	Type    string      `json:"type"`    // message, join, leave, system
	RoomID  string      `json:"room_id"`
	UserID  string      `json:"user_id"`
	Data    interface{} `json:"data"`
}

// NewWSManager 创建 WebSocket 管理器
func NewWSManager() *WSManager {
	return &WSManager{
		clients:    make(map[string]*WSClient),
		rooms:      make(map[string]map[string]*WSClient),
		broadcast:  make(chan *WSMessage, 256),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
}

// Run 启动 WebSocket 管理器
func (m *WSManager) Run() {
	for {
		select {
		case client := <-m.register:
			m.handleRegister(client)

		case client := <-m.unregister:
			m.handleUnregister(client)

		case message := <-m.broadcast:
			m.handleBroadcast(message)
		}
	}
}

// handleRegister 处理客户端注册
func (m *WSManager) handleRegister(client *WSClient) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 注册到全局客户端
	m.clients[client.userID] = client

	// 注册到房间
	if client.roomID != "" {
		if m.rooms[client.roomID] == nil {
			m.rooms[client.roomID] = make(map[string]*WSClient)
		}
		m.rooms[client.roomID][client.userID] = client
	}

	log.Printf("WebSocket client registered: user=%s, room=%s", client.userID, client.roomID)
}

// handleUnregister 处理客户端注销
func (m *WSManager) handleUnregister(client *WSClient) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 从全局客户端移除
	if _, ok := m.clients[client.userID]; ok {
		delete(m.clients, client.userID)
		close(client.send)
	}

	// 从房间移除
	if client.roomID != "" {
		if room, ok := m.rooms[client.roomID]; ok {
			delete(room, client.userID)
			if len(room) == 0 {
				delete(m.rooms, client.roomID)
			}
		}
	}

	log.Printf("WebSocket client unregistered: user=%s", client.userID)
}

// handleBroadcast 处理广播消息
func (m *WSManager) handleBroadcast(message *WSMessage) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	// 如果是房间消息，只广播给房间内的客户端
	if message.RoomID != "" {
		if room, ok := m.rooms[message.RoomID]; ok {
			for _, client := range room {
				select {
				case client.send <- data:
				default:
					// 客户端发送缓冲区满，关闭连接
					close(client.send)
				}
			}
		}
	}
}

// BroadcastToRoom 向房间广播消息
func (m *WSManager) BroadcastToRoom(roomID string, msgType string, data interface{}) {
	message := &WSMessage{
		Type:   msgType,
		RoomID: roomID,
		Data:   data,
	}
	m.broadcast <- message
}

// GetOnlineCount 获取房间在线人数
func (m *WSManager) GetOnlineCount(roomID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if room, ok := m.rooms[roomID]; ok {
		return len(room)
	}
	return 0
}

// GetOnlineUsers 获取房间在线用户列表
func (m *WSManager) GetOnlineUsers(roomID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var users []string
	if room, ok := m.rooms[roomID]; ok {
		for userID := range room {
			users = append(users, userID)
		}
	}
	return users
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域
	},
}

// HandleWebSocket WebSocket 连接处理
func (m *WSManager) HandleWebSocket(w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &WSClient{
		manager: m,
		conn:    conn,
		send:    make(chan []byte, 256),
		userID:  userID,
	}

	m.register <- client

	// 启动读写 goroutine
	go client.writePump()
	go client.readPump()
}

// readPump 读取客户端消息
func (c *WSClient) readPump() {
	defer func() {
		c.manager.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024) // 512KB

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// 解析消息
		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// 设置发送者
		msg.UserID = c.userID
		c.roomID = msg.RoomID

		// 广播消息
		c.manager.broadcast <- &msg
	}
}

// writePump 向客户端写入消息
func (c *WSClient) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// Global WSManager instance
var GlobalWSManager = NewWSManager()

func init() {
	go GlobalWSManager.Run()
}
