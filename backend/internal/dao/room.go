package dao

import (
	"github.com/baijianruoli/bot_chat/backend/internal/model"
	"gorm.io/gorm"
)

// RoomDAO 房间数据访问对象
type RoomDAO struct {
	db *gorm.DB
}

// NewRoomDAO 创建 RoomDAO
func NewRoomDAO(db *gorm.DB) *RoomDAO {
	return &RoomDAO{db: db}
}

// Create 创建房间
func (d *RoomDAO) Create(room *model.Room) error {
	return d.db.Create(room).Error
}

// GetByID 根据ID获取房间
func (d *RoomDAO) GetByID(roomID string) (*model.Room, error) {
	var room model.Room
	err := d.db.Where("room_id = ?", roomID).First(&room).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &room, err
}

// List 获取房间列表
func (d *RoomDAO) List(page, pageSize int) ([]*model.Room, int64, error) {
	var rooms []*model.Room
	var total int64
	
	offset := (page - 1) * pageSize
	
	err := d.db.Model(&model.Room{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	err = d.db.Offset(offset).Limit(pageSize).Find(&rooms).Error
	return rooms, total, err
}

// UpdateUserCount 更新房间人数
func (d *RoomDAO) UpdateUserCount(roomID string, delta int32) error {
	return d.db.Model(&model.Room{}).
		Where("room_id = ?", roomID).
		Update("user_count", gorm.Expr("user_count + ?", delta)).Error
}

// RoomMemberDAO 房间成员 DAO
type RoomMemberDAO struct {
	db *gorm.DB
}

// NewRoomMemberDAO 创建 RoomMemberDAO
func NewRoomMemberDAO(db *gorm.DB) *RoomMemberDAO {
	return &RoomMemberDAO{db: db}
}

// AddMember 添加成员
func (d *RoomMemberDAO) AddMember(roomID, userID string) error {
	member := &model.RoomMember{
		RoomID: roomID,
		UserID: userID,
	}
	return d.db.Create(member).Error
}

// RemoveMember 移除成员
func (d *RoomMemberDAO) RemoveMember(roomID, userID string) error {
	return d.db.Where("room_id = ? AND user_id = ?", roomID, userID).
		Delete(&model.RoomMember{}).Error
}

// IsMember 检查是否是成员
func (d *RoomMemberDAO) IsMember(roomID, userID string) (bool, error) {
	var count int64
	err := d.db.Model(&model.RoomMember{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetMembers 获取房间成员列表
func (d *RoomMemberDAO) GetMembers(roomID string) ([]string, error) {
	var userIDs []string
	err := d.db.Model(&model.RoomMember{}).
		Where("room_id = ?", roomID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
