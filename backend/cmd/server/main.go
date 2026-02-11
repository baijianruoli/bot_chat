package main

import (
	"log"

	"github.com/baijianruoli/bot_chat/backend/internal/conf"
)

func main() {
	// 加载配置
	config := conf.LoadConfig()
	log.Printf("Bot Chat Server starting on %s:%d", config.Server.Host, config.Server.Port)
	
	// TODO: 初始化Kitex服务
	// TODO: 注册Handler
	// TODO: 启动服务
	
	log.Println("Server started successfully!")
	
	// 阻塞主线程
	select {}
}
