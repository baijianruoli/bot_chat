package conf

import (
	"os"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port int
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Config 全局配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 从环境变量读取，或使用默认值
	GlobalConfig = &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvAsInt("SERVER_PORT", 8888),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			Username: getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASS", ""),
			Database: getEnv("DB_NAME", "bot_chat"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASS", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}
	return GlobalConfig
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	// 简化处理，实际应转换
	return defaultVal
}
