package main

import (
	"log"
	"net"
	"net/http"
	
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
	db, err := dao.InitDB()
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	log.Println("Database initialized")
	
	// 创建服务实例
	svc := service.NewChatService()
	
	// 启动 WebSocket 管理器
	go service.GlobalWSManager.Run()
	log.Println("WebSocket manager started")
	
	// 启动 HTTP 服务器（用于 WebSocket 和 API）
	go startHTTPServer(config.Server.Host, config.Server.Port + 1)
	
	// 创建 Kitex gRPC 服务
	svr := server.NewServer(
		server.WithServiceAddr(&net.TCPAddr{
			IP:   net.ParseIP(config.Server.Host),
			Port: config.Server.Port,
		}),
	)
	
	// 注册服务
	chat.RegisterService(svr, svc)
	
	log.Println("Server started successfully!")
	log.Printf("gRPC server listening on port %d", config.Server.Port)
	log.Printf("HTTP/WebSocket server listening on port %d", config.Server.Port + 1)
	
	// 启动服务
	if err := svr.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}

// startHTTPServer 启动 HTTP 服务器
func startHTTPServer(host string, port int) {
	// WebSocket 路由
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		service.GlobalWSManager.HandleWebSocket(w, r, r.URL.Query().Get("user_id"))
	})
	
	// 健康检查
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	addr := net.JoinHostPort(host, string(rune('0' + port - 8888)))
	log.Printf("HTTP server starting on %s", addr)
	
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
