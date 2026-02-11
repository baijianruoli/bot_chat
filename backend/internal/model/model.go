package model

// User 用户模型
type User struct {
	UserID    string `json:"user_id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"uniqueIndex;not null"`
	Password  string `json:"-" gorm:"not null"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt int64  `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

// Room 聊天室模型
type Room struct {
	RoomID      string `json:"room_id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	CreatorID   string `json:"creator_id"`
	UserCount   int32  `json:"user_count" gorm:"default:0"`
	CreatedAt   int64  `json:"created_at" gorm:"autoCreateTime:milli"`
	UpdatedAt   int64  `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

// RoomMember 房间成员关系
type RoomMember struct {
	ID        uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
	RoomID    string `json:"room_id" gorm:"index"`
	UserID    string `json:"user_id" gorm:"index"`
	JoinTime  int64  `json:"join_time" gorm:"autoCreateTime:milli"`
}

// Message 消息模型
type Message struct {
	MsgID     string `json:"msg_id" gorm:"primaryKey"`
	RoomID    string `json:"room_id" gorm:"index"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	MsgType   int32  `json:"msg_type" gorm:"default:1"` // 1:文本 2:图片 3:系统
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime:milli"`
	
	// 关联用户（不存入数据库）
	User *User `json:"user,omitempty" gorm:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func (Room) TableName() string {
	return "rooms"
}

func (RoomMember) TableName() string {
	return "room_members"
}

func (Message) TableName() string {
	return "messages"
}
