package main

import (
	"log"
	"net"
	
	"github.com/baijianruoli/bot_chat/backend/internal/conf"
	"github.com/baijianruoli/bot_chat/backend/internal/dao"
	"github.com/baijianruoli/bot_chat/backend/internal/service"
	chat "github.com/baijianruoli/bot_chat/backend/kitex_gen/chat"
	"github.com/cloudwego/kitex/server"
)

func main() {
	// 加载配置
	config := conf.LoadConfig()
	log.Printf("Bot Chat Server starting on %s:%d", config.Server.Host, config.Server.Port)
	
	// 初始化数据库
	_, err := dao.InitDB()
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	log.Println("Database initialized")
	
	// 创建服务实例
	svc := service.NewChatService()
	
	// 创建 Kitex 服务
	svr := server.NewServer(
		server.WithServiceAddr(&net.TCPAddr{
			IP:   net.ParseIP(config.Server.Host),
			Port: config.Server.Port,
		}),
	)
	
	// 注册服务
	chat.RegisterService(svr, svc)
	
	log.Println("Server started successfully!")
	
	// 启动服务
	if err := svr.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
