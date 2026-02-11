package service

import (
	"context"
	"errors"
	
	"github.com/baijianruoli/bot_chat/backend/internal/dao"
	"github.com/baijianruoli/bot_chat/backend/internal/model"
	"github.com/baijianruoli/bot_chat/backend/internal/utils"
	chat "github.com/baijianruoli/bot_chat/backend/kitex_gen/chat"
)

// ChatServiceImpl 聊天服务实现
type ChatServiceImpl struct{}

// NewChatService 创建服务实例
func NewChatService() *ChatServiceImpl {
	return &ChatServiceImpl{}
}

// Register 用户注册
func (s *ChatServiceImpl) Register(ctx context.Context, req *chat.RegisterReq) (*chat.RegisterResp, error) {
	// 检查用户名是否已存在
	userDAO := dao.NewUserDAO(dao.DB)
	existingUser, err := userDAO.GetByUsername(req.Username)
	if err != nil {
		return utils.Error(utils.CodeServerError, "database error").(*utils.Resp), nil
	}
	if existingUser != nil {
		return &chat.RegisterResp{
			Code:    utils.CodeUserExists,
			Message: "username already exists",
		}, nil
	}
	
	// 创建新用户
	user := &model.User{
		UserID:   utils.GenerateUserID(),
		Username: req.Username,
		Password: utils.HashPassword(req.Password),
		Nickname: req.Nickname,
	}
	
	if err := userDAO.Create(user); err != nil {
		return &chat.RegisterResp{
			Code:    utils.CodeServerError,
			Message: "failed to create user",
		}, nil
	}
	
	return &chat.RegisterResp{
		Code:    utils.CodeSuccess,
		Message: "success",
		UserId:  user.UserID,
	}, nil
}

// Login 用户登录
func (s *ChatServiceImpl) Login(ctx context.Context, req *chat.LoginReq) (*chat.LoginResp, error) {
	userDAO := dao.NewUserDAO(dao.DB)
	user, err := userDAO.GetByUsername(req.Username)
	if err != nil {
		return &chat.LoginResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	if user == nil {
		return &chat.LoginResp{
			Code:    utils.CodeUserNotFound,
			Message: "user not found",
		}, nil
	}
	
	// 验证密码
	if !utils.VerifyPassword(req.Password, user.Password) {
		return &chat.LoginResp{
			Code:    utils.CodePasswordError,
			Message: "password incorrect",
		}, nil
	}
	
	// 生成 token
	token := utils.MD5(user.UserID + utils.GetCurrentTimeStr())
	
	return &chat.LoginResp{
		Code:    utils.CodeSuccess,
		Message: "success",
		Token:   token,
		User: &chat.UserInfo{
			UserId:    user.UserID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// CreateRoom 创建房间
func (s *ChatServiceImpl) CreateRoom(ctx context.Context, req *chat.CreateRoomReq) (*chat.CreateRoomResp, error) {
	roomDAO := dao.NewRoomDAO(dao.DB)
	roomMemberDAO := dao.NewRoomMemberDAO(dao.DB)
	
	room := &model.Room{
		RoomID:      utils.GenerateRoomID(),
		Name:        req.Name,
		Description: req.Description,
		CreatorID:   req.CreatorId,
		UserCount:   1,
	}
	
	if err := roomDAO.Create(room); err != nil {
		return &chat.CreateRoomResp{
			Code:    utils.CodeServerError,
			Message: "failed to create room",
		}, nil
	}
	
	// 创建者自动加入房间
	if err := roomMemberDAO.AddMember(room.RoomID, req.CreatorId); err != nil {
		return &chat.CreateRoomResp{
			Code:    utils.CodeServerError,
			Message: "failed to join room",
		}, nil
	}
	
	return &chat.CreateRoomResp{
		Code:    utils.CodeSuccess,
		Message: "success",
		Room: &chat.RoomInfo{
			RoomId:    room.RoomID,
			Name:      room.Name,
			Description: room.Description,
			CreatorId: room.CreatorID,
			UserCount: room.UserCount,
			CreatedAt: room.CreatedAt,
		},
	}, nil
}

// ListRooms 获取房间列表
func (s *ChatServiceImpl) ListRooms(ctx context.Context, req *chat.ListRoomsReq) (*chat.ListRoomsResp, error) {
	roomDAO := dao.NewRoomDAO(dao.DB)
	
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	
	rooms, total, err := roomDAO.List(int(req.Page), int(req.PageSize))
	if err != nil {
		return &chat.ListRoomsResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	
	roomList := make([]*chat.RoomInfo, len(rooms))
	for i, room := range rooms {
		roomList[i] = &chat.RoomInfo{
			RoomId:      room.RoomID,
			Name:        room.Name,
			Description: room.Description,
			CreatorId:   room.CreatorID,
			UserCount:   room.UserCount,
			CreatedAt:   room.CreatedAt,
		}
	}
	
	return &chat.ListRoomsResp{
		Code:    utils.CodeSuccess,
		Message: "success",
		Rooms:   roomList,
		Total:   int32(total),
	}, nil
}

// JoinRoom 加入房间
func (s *ChatServiceImpl) JoinRoom(ctx context.Context, req *chat.JoinRoomReq) (*chat.JoinRoomResp, error) {
	roomDAO := dao.NewRoomDAO(dao.DB)
	roomMemberDAO := dao.NewRoomMemberDAO(dao.DB)
	
	// 检查房间是否存在
	room, err := roomDAO.GetByID(req.RoomId)
	if err != nil {
		return &chat.JoinRoomResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	if room == nil {
		return &chat.JoinRoomResp{
			Code:    utils.CodeRoomNotFound,
			Message: "room not found",
		}, nil
	}
	
	// 检查是否已在房间中
	isMember, err := roomMemberDAO.IsMember(req.RoomId, req.UserId)
	if err != nil {
		return &chat.JoinRoomResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	if isMember {
		return &chat.JoinRoomResp{
			Code:    utils.CodeAlreadyInRoom,
			Message: "already in room",
		}, nil
	}
	
	// 添加成员
	if err := roomMemberDAO.AddMember(req.RoomId, req.UserId); err != nil {
		return &chat.JoinRoomResp{
			Code:    utils.CodeServerError,
			Message: "failed to join room",
		}, nil
	}
	
	// 更新房间人数
	roomDAO.UpdateUserCount(req.RoomId, 1)
	
	return &chat.JoinRoomResp{
		Code:    utils.CodeSuccess,
		Message: "success",
		Room: &chat.RoomInfo{
			RoomId:      room.RoomID,
			Name:        room.Name,
			Description: room.Description,
			CreatorId:   room.CreatorID,
			UserCount:   room.UserCount + 1,
			CreatedAt:   room.CreatedAt,
		},
	}, nil
}

// LeaveRoom 离开房间
func (s *ChatServiceImpl) LeaveRoom(ctx context.Context, req *chat.LeaveRoomReq) (*chat.LeaveRoomResp, error) {
	roomDAO := dao.NewRoomDAO(dao.DB)
	roomMemberDAO := dao.NewRoomMemberDAO(dao.DB)
	
	// 检查是否在房间中
	isMember, err := roomMemberDAO.IsMember(req.RoomId, req.UserId)
	if err != nil {
		return &chat.LeaveRoomResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	if !isMember {
		return &chat.LeaveRoomResp{
			Code:    utils.CodeNotInRoom,
			Message: "not in room",
		}, nil
	}
	
	// 移除成员
	if err := roomMemberDAO.RemoveMember(req.RoomId, req.UserId); err != nil {
		return &chat.LeaveRoomResp{
			Code:    utils.CodeServerError,
			Message: "failed to leave room",
		}, nil
	}
	
	// 更新房间人数
	roomDAO.UpdateUserCount(req.RoomId, -1)
	
	return &chat.LeaveRoomResp{
		Code:    utils.CodeSuccess,
		Message: "success",
	}, nil
}

// SendMessage 发送消息
func (s *ChatServiceImpl) SendMessage(ctx context.Context, req *chat.SendMessageReq) (*chat.SendMessageResp, error) {
	roomDAO := dao.NewRoomDAO(dao.DB)
	roomMemberDAO := dao.NewRoomMemberDAO(dao.DB)
	userDAO := dao.NewUserDAO(dao.DB)
	messageDAO := dao.NewMessageDAO(dao.DB)
	
	// 检查房间是否存在
	room, err := roomDAO.GetByID(req.RoomId)
	if err != nil {
		return &chat.SendMessageResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	if room == nil {
		return &chat.SendMessageResp{
			Code:    utils.CodeRoomNotFound,
			Message: "room not found",
		}, nil
	}
	
	// 检查用户是否在房间中
	isMember, err := roomMemberDAO.IsMember(req.RoomId, req.UserId)
	if err != nil {
		return &chat.SendMessageResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	if !isMember {
		return &chat.SendMessageResp{
			Code:    utils.CodeNotInRoom,
			Message: "not in room",
		}, nil
	}
	
	// 获取发送者信息
	user, err := userDAO.GetByID(req.UserId)
	if err != nil {
		return &chat.SendMessageResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	
	// 创建消息
	msg := &model.Message{
		MsgID:   utils.GenerateMsgID(),
		RoomID:  req.RoomId,
		UserID:  req.UserId,
		Content: req.Content,
		MsgType: req.MsgType,
	}
	
	if err := messageDAO.Create(msg); err != nil {
		return &chat.SendMessageResp{
			Code:    utils.CodeServerError,
			Message: "failed to send message",
		}, nil
	}
	
	return &chat.SendMessageResp{
		Code:    utils.CodeSuccess,
		Message: "success",
		Msg: &chat.MessageInfo{
			MsgId:     msg.MsgID,
			RoomId:    msg.RoomID,
			Sender: &chat.UserInfo{
				UserId:   user.UserID,
				Username: user.Username,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
			},
			Content:   msg.Content,
			MsgType:   msg.MsgType,
			Timestamp: msg.CreatedAt,
		},
	}, nil
}

// GetHistory 获取历史消息
func (s *ChatServiceImpl) GetHistory(ctx context.Context, req *chat.GetHistoryReq) (*chat.GetHistoryResp, error) {
	roomDAO := dao.NewRoomDAO(dao.DB)
	roomMemberDAO := dao.NewRoomMemberDAO(dao.DB)
	userDAO := dao.NewUserDAO(dao.DB)
	messageDAO := dao.NewMessageDAO(dao.DB)
	
	// 检查房间是否存在
	room, err := roomDAO.GetByID(req.RoomId)
	if err != nil {
		return &chat.GetHistoryResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	if room == nil {
		return &chat.GetHistoryResp{
			Code:    utils.CodeRoomNotFound,
			Message: "room not found",
		}, nil
	}
	
	// 检查用户是否在房间中
	isMember, err := roomMemberDAO.IsMember(req.RoomId, req.UserId)
	if err != nil {
		return &chat.GetHistoryResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	if !isMember {
		return &chat.GetHistoryResp{
			Code:    utils.CodeNotInRoom,
			Message: "not in room",
		}, nil
	}
	
	// 获取消息
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	
	messages, err := messageDAO.GetHistory(req.RoomId, req.BeforeTime, int(req.Limit))
	if err != nil {
		return &chat.GetHistoryResp{
			Code:    utils.CodeServerError,
			Message: "database error",
		}, nil
	}
	
	// 填充发送者信息
	msgList := make([]*chat.MessageInfo, len(messages))
	for i, msg := range messages {
		user, _ := userDAO.GetByID(msg.UserID)
		msgList[i] = &chat.MessageInfo{
			MsgId:  msg.MsgID,
			RoomId: msg.RoomID,
			Sender: &chat.UserInfo{
				UserId:   user.UserID,
				Username: user.Username,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
			},
			Content:   msg.Content,
			MsgType:   msg.MsgType,
			Timestamp: msg.CreatedAt,
		}
	}
	
	hasMore := len(messages) == int(req.Limit)
	
	return &chat.GetHistoryResp{
		Code:     utils.CodeSuccess,
		Message:  "success",
		Messages: msgList,
		HasMore:  hasMore,
	}, nil
}
