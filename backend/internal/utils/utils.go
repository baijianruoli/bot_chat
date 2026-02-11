package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID 生成唯一ID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateUserID 生成用户ID
func GenerateUserID() string {
	return "u_" + GenerateUUID()[:8]
}

// GenerateRoomID 生成房间ID
func GenerateRoomID() string {
	return "r_" + GenerateUUID()[:8]
}

// GenerateMsgID 生成消息ID
func GenerateMsgID() string {
	return "m_" + GenerateUUID()[:12]
}

// MD5 计算MD5哈希
func MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// HashPassword 密码哈希
func HashPassword(password string) string {
	salt := "bot_chat_salt_2024"
	return MD5(password + salt)
}

// VerifyPassword 验证密码
func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// GetCurrentTimestamp 获取当前时间戳（毫秒）
func GetCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// FormatTime 格式化时间
func FormatTime(timestamp int64) string {
	t := time.Unix(timestamp/1000, 0)
	return t.Format("2006-01-02 15:04:05")
}

// GetCurrentTimeStr 获取当前时间字符串
func GetCurrentTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Resp 通用响应结构
type Resp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) *Resp {
	return &Resp{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// Error 错误响应
func Error(code int32, message string) *Resp {
	return &Resp{
		Code:    code,
		Message: message,
	}
}

// Errorf 格式化错误响应
func Errorf(code int32, format string, args ...interface{}) *Resp {
	return &Resp{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Common errors
const (
	CodeSuccess        = 0
	CodeParamError     = 400
	CodeUnauthorized   = 401
	CodeNotFound       = 404
	CodeServerError    = 500
	CodeUserExists     = 1001
	CodeUserNotFound   = 1002
	CodePasswordError  = 1003
	CodeRoomNotFound   = 2001
	CodeRoomExists     = 2002
	CodeAlreadyInRoom  = 2003
	CodeNotInRoom      = 2004
)
