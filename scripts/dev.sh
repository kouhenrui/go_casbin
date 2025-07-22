#!/bin/bash

# 设置开发环境变量
export GO_ENV=development
export DEBUG=true

# 检查 air 是否安装
if ! command -v air &> /dev/null; then
    echo "air 未安装，正在安装..."
    go install github.com/air-verse/air@latest
fi

# 创建临时目录
mkdir -p tmp

# 启动 air 热重载
echo "启动开发模式热重载..."
air 