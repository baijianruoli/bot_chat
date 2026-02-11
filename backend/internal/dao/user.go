package dao

import (
	"github.com/baijianruoli/bot_chat/backend/internal/model"
	"gorm.io/gorm"
)

// UserDAO 用户数据访问对象
type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO 创建 UserDAO
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

// Create 创建用户
func (d *UserDAO) Create(user *model.User) error {
	return d.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (d *UserDAO) GetByID(userID string) (*model.User, error) {
	var user model.User
	err := d.db.Where("user_id = ?", userID).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// GetByUsername 根据用户名获取用户
func (d *UserDAO) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := d.db.Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

// Update 更新用户信息
func (d *UserDAO) Update(user *model.User) error {
	return d.db.Save(user).Error
}
