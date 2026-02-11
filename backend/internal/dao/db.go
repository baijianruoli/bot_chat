package dao

import (
	"fmt"
	"log"
	
	"github.com/baijianruoli/bot_chat/backend/internal/conf"
	"github.com/baijianruoli/bot_chat/backend/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() (*gorm.DB, error) {
	config := conf.GlobalConfig.Database
	
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}
	
	// 自动迁移
	err = db.AutoMigrate(
		&model.User{},
		&model.Room{},
		&model.RoomMember{},
		&model.Message{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}
	
	DB = db
	log.Println("Database connected successfully")
	return db, nil
}
