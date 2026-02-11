# Bot Chat

一个基于 **Golang + Kitex + React** 的实时聊天室项目。

## 🏗️ 项目架构

```
bot_chat/
├── backend/          # 后端服务 (Golang + Kitex)
│   ├── cmd/          # 入口程序
│   ├── internal/     # 内部实现
│   ├── kitex_gen/    # Kitex生成的代码
│   └── proto/        # Protocol Buffers定义
├── frontend/         # 前端应用 (React)
│   ├── src/          # 源代码
│   ├── public/       # 静态资源
│   └── package.json  # 依赖配置
└── docker-compose.yml # 部署配置
```

## 🚀 技术栈

### 后端
- **语言**: Golang 1.21+
- **RPC框架**: [Kitex](https://github.com/cloudwego/kitex) (字节跳动开源)
- **通信**: WebSocket + gRPC
- **存储**: Redis (消息缓存) + MySQL (用户数据)
- **消息队列**: Kafka (可选)

### 前端
- **框架**: React 18
- **构建工具**: Vite
- **UI组件**: Ant Design
- **状态管理**: Zustand
- **WebSocket**: Socket.io-client

## 📦 功能特性

- [x] 用户注册/登录
- [x] 实时消息收发
- [x] 多房间支持
- [x] 在线用户列表
- [x] 消息历史记录
- [x] 心跳保活机制

## 🛠️ 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- Redis 6+
- MySQL 8+

### 后端启动
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 前端启动
```bash
cd frontend
npm install
npm run dev
```

## 📄 许可证

MIT License
