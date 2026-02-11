package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/baijianruoli/bot_chat/backend/internal/dao"
	"github.com/baijianruoli/bot_chat/backend/internal/model"
	"github.com/baijianruoli/bot_chat/backend/internal/utils"
	chat "github.com/baijianruoli/bot_chat/backend/kitex_gen/chat"
)

// SendMessageWithWS 发送消息并通过 WebSocket 广播
func (s *ChatServiceImpl) SendMessageWithWS(ctx context.Context, req *chat.SendMessageReq) (*chat.SendMessageResp, error) {
	// 先调用原有的 SendMessage 逻辑
	resp, err := s.SendMessage(ctx, req)
	if err != nil || resp.Code != utils.CodeSuccess {
		return resp, err
	}

	// 通过 WebSocket 广播消息
	if resp.Msg != nil {
		message := &WSMessage{
			Type:   "message",
			RoomID: req.RoomId,
			UserID: req.UserId,
			Data:   resp.Msg,
		}
		GlobalWSManager.broadcast <- message
	}

	return resp, nil
}

// JoinRoomWithWS 加入房间并通过 WebSocket 广播
func (s *ChatServiceImpl) JoinRoomWithWS(ctx context.Context, req *chat.JoinRoomReq, wsClient *WSClient) (*chat.JoinRoomResp, error) {
	resp, err := s.JoinRoom(ctx, req)
	if err != nil || resp.Code != utils.CodeSuccess {
		return resp, err
	}

	// 更新 WebSocket 客户端的房间
	wsClient.roomID = req.RoomId

	// 广播用户加入消息
	GlobalWSManager.BroadcastToRoom(req.RoomId, "join", map[string]interface{}{
		"user_id":  req.UserId,
		"nickname": resp.Room.Name,
	})

	// 更新在线人数
	GlobalWSManager.BroadcastToRoom(req.RoomId, "online_count", map[string]interface{}{
		"count": GlobalWSManager.GetOnlineCount(req.RoomId),
	})

	return resp, nil
}

// LeaveRoomWithWS 离开房间并通过 WebSocket 广播
func (s *ChatServiceImpl) LeaveRoomWithWS(ctx context.Context, req *chat.LeaveRoomReq) (*chat.LeaveRoomResp, error) {
	resp, err := s.LeaveRoom(ctx, req)
	if err != nil || resp.Code != utils.CodeSuccess {
		return resp, err
	}

	// 广播用户离开消息
	GlobalWSManager.BroadcastToRoom(req.RoomId, "leave", map[string]interface{}{
		"user_id": req.UserId,
	})

	// 更新在线人数
	GlobalWSManager.BroadcastToRoom(req.RoomId, "online_count", map[string]interface{}{
		"count": GlobalWSManager.GetOnlineCount(req.RoomId),
	})

	return resp, nil
}

// WSMessageData WebSocket 消息数据结构
type WSMessageData struct {
	MsgID     string `json:"msg_id"`
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Timestamp int64  `json:"timestamp"`
	MsgType   int32  `json:"msg_type"`
}

// HandleWSMessage 处理 WebSocket 消息
func HandleWSMessage(userID string, roomID string, data json.RawMessage) {
	var msgData WSMessageData
	if err := json.Unmarshal(data, &msgData); err != nil {
		return
	}

	// 保存消息到数据库
	msg := &model.Message{
		MsgID:   utils.GenerateMsgID(),
		RoomID:  roomID,
		UserID:  userID,
		Content: msgData.Content,
		MsgType: msgData.MsgType,
	}

	dao := dao.NewMessageDAO(dao.DB)
	dao.Create(msg)

	// 广播消息
	GlobalWSManager.BroadcastToRoom(roomID, "message", msg)
}

// WSRouter WebSocket 路由
type WSRouter struct {
	chatService *ChatServiceImpl
}

// NewWSRouter 创建 WebSocket 路由
func NewWSRouter() *WSRouter {
	return &WSRouter{
		chatService: NewChatService(),
	}
}

// HandleConnection 处理 WebSocket 连接
func (r *WSRouter) HandleConnection(w http.ResponseWriter, req *http.Request) {
	// 从 query 参数获取 userID
	userID := req.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id required", http.StatusBadRequest)
		return
	}

	// 升级为 WebSocket
	GlobalWSManager.HandleWebSocket(w, req, userID)
}
