.PHONY: all backend frontend docker clean

# 默认目标
all: backend frontend

# 后端构建
backend:
	@echo "Building backend..."
	cd backend && go mod tidy && go build -o bin/server cmd/server/main.go

# 前端构建
frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

# 生成 Kitex 代码
generate:
	@echo "Generating Kitex code..."
	cd backend && kitex -module github.com/baijianruoli/bot_chat/backend \
		-service chat_service \
		-type protobuf \
		proto/chat.proto

# Docker 构建
docker:
	@echo "Building Docker images..."
	docker-compose build

# Docker 运行
docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

# Docker 停止
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

# 清理
clean:
	@echo "Cleaning..."
	rm -rf backend/bin
	rm -rf frontend/dist
	rm -rf frontend/node_modules

# 运行后端un-backend:
	@echo "Running backend..."
	cd backend && go run cmd/server/main.go

# 运行前端un-frontend:
	@echo "Running frontend..."
	cd frontend && npm run dev

# 测试
test:
	@echo "Running tests..."
	cd backend && go test ./...
	cd frontend && npm test
