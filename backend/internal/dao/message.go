package dao

import (
	"github.com/baijianruoli/bot_chat/backend/internal/model"
)

// MessageDAO 消息数据访问对象
type MessageDAO struct {
	db *gorm.DB
}

// NewMessageDAO 创建 MessageDAO
func NewMessageDAO(db *gorm.DB) *MessageDAO {
	return &MessageDAO{db: db}
}

// Create 创建消息
func (d *MessageDAO) Create(msg *model.Message) error {
	return d.db.Create(msg).Error
}

// GetHistory 获取历史消息
func (d *MessageDAO) GetHistory(roomID string, beforeTime int64, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	
	query := d.db.Where("room_id = ?", roomID)
	
	if beforeTime > 0 {
		query = query.Where("created_at < ?", beforeTime)
	}
	
	err := query.Order("created_at DESC").Limit(limit).Find(&messages).Error
	
	// 反转顺序
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	
	return messages, err
}

// GetByID 根据ID获取消息
func (d *MessageDAO) GetByID(msgID string) (*model.Message, error) {
	var msg model.Message
	err := d.db.Where("msg_id = ?", msgID).First(&msg).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &msg, err
}
