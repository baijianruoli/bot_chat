#!/bin/bash
# 生成 Kitex 代码

cd backend

echo "Generating Kitex code from proto..."

# 安装 kitex 工具（如果未安装）
if ! command -v kitex &> /dev/null; then
    echo "Installing kitex..."
    go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
fi

# 生成代码
kitex -module github.com/baijianruoli/bot_chat/backend \
      -service chat_service \
      -type protobuf \
      proto/chat.proto

echo "Done!"
