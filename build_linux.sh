#!/bin/bash

echo "=========================================="
echo "Start compiling Authos for Linux (AMD64)..."
echo "=========================================="

# 设置 Go 环境变量
# CGO_ENABLED=0: 禁用 CGO，使用纯 Go 实现
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

# 执行编译
# -ldflags="-s -w": 去除符号表和调试信息，减小体积
# -trimpath: 移除构建路径信息
go build -ldflags="-s -w" -trimpath -o authos-linux-amd64 .

if [ $? -ne 0 ]; then
    echo "[ERROR] Compilation failed!"
    exit 1
fi

echo "=========================================="
echo "Compilation successful!"
echo "Output file: authos-linux-amd64"
echo "=========================================="
