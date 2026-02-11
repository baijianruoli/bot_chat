# Bot Chat

一个基于 **Golang + Kitex + React** 的实时聊天室项目。

## 🏗️ 项目架构

```
bot_chat/
├── backend/          # 后端服务 (Golang + Kitex)
│   ├── cmd/          # 入口程序
│   ├── internal/     # 内部实现
│   │   ├── conf/     # 配置
│   │   ├── dao/      # 数据访问层
│   │   ├── model/    # 数据模型
│   │   ├── service/  # 业务逻辑层
│   │   └── utils/    # 工具函数
│   ├── kitex_gen/    # Kitex生成的代码
│   └── proto/        # Protocol Buffers定义
├── frontend/         # 前端应用 (React)
│   ├── src/          # 源代码
│   │   ├── api/      # API接口
│   │   ├── components/# 组件
│   │   ├── pages/    # 页面
│   │   └── store/    # 状态管理
│   └── public/       # 静态资源
└── docker-compose.yml # 部署配置
```

## 🚀 技术栈

### 后端
- **语言**: Golang 1.21+
- **RPC框架**: [Kitex](https://github.com/cloudwego/kitex) (字节跳动开源)
- **通信**: gRPC + HTTP
- **存储**: MySQL (用户数据) + Redis (缓存)
- **ORM**: GORM

### 前端
- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **UI组件**: Ant Design
- **状态管理**: Zustand
- **HTTP客户端**: Axios
- **日期处理**: Dayjs

## 📦 功能特性

- ✅ 用户注册/登录
- ✅ JWT Token 认证
- ✅ 创建/加入聊天室
- ✅ 实时消息收发
- ✅ 消息历史记录
- ✅ 在线用户列表
- ✅ 响应式设计

## 🛠️ 快速开始

### 方式1：Docker Compose（推荐）

```bash
# 克隆项目
git clone https://github.com/baijianruoli/bot_chat.git
cd bot_chat

# 启动所有服务
docker-compose up -d

# 访问
# 前端: http://localhost:3000
# 后端: http://localhost:8888
```

### 方式2：本地开发

#### 环境要求
- Go 1.21+
- Node.js 18+
- MySQL 8+
- Redis 6+

#### 后端启动

```bash
cd backend

# 安装依赖
go mod tidy

# 生成 Kitex 代码（如果需要）
chmod +x ../scripts/generate-kitex.sh
../scripts/generate-kitex.sh

# 配置数据库
# 修改 internal/conf/conf.go 中的数据库配置

# 运行
go run cmd/server/main.go
```

#### 前端启动

```bash
cd frontend

# 安装依赖
npm install

# 开发模式
npm run dev

# 构建
npm run build
```

## 📡 API 接口

### 用户相关
- `POST /register` - 用户注册
- `POST /login` - 用户登录

### 房间相关
- `GET /rooms` - 获取房间列表
- `POST /rooms` - 创建房间
- `POST /rooms/:id/join` - 加入房间
- `POST /rooms/:id/leave` - 离开房间

### 消息相关
- `GET /messages?room_id=xxx` - 获取历史消息
- `POST /messages` - 发送消息

## 📁 项目结构说明

```
backend/
├── cmd/server/main.go      # 服务入口
├── internal/
│   ├── conf/conf.go        # 配置管理
│   ├── dao/                # 数据访问层
│   │   ├── db.go          # 数据库连接
│   │   ├── user.go        # 用户DAO
│   │   ├── room.go        # 房间DAO
│   │   └── message.go     # 消息DAO
│   ├── model/model.go     # 数据模型
│   ├── service/chat.go    # 业务逻辑实现
│   └── utils/utils.go     # 工具函数
└── proto/chat.proto       # gRPC 协议定义

frontend/
├── src/
│   ├── api/index.ts       # API 接口封装
│   ├── components/        # 公共组件
│   │   └── Layout.tsx     # 布局组件
│   ├── pages/             # 页面
│   │   ├── Login.tsx      # 登录/注册
│   │   ├── RoomList.tsx   # 房间列表
│   │   └── Chat.tsx       # 聊天页面
│   ├── store/index.ts     # Zustand 状态管理
│   └── App.tsx            # 应用入口
```

## 🔧 开发计划

- [x] 项目基础架构
- [x] 后端基础实现
- [x] 前端基础实现
- [ ] WebSocket 实时通信
- [ ] 消息已读状态
- [ ] 文件上传
- [ ] 用户头像上传
- [ ] 消息撤回
- [ ] 私聊功能

## 📄 许可证

MIT License
