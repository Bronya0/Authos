#!/bin/bash

# 应用名称（即编译后的二进制文件名）
APP_NAME="authos-linux-amd64"

# 切换到脚本所在目录，确保相对路径正确
cd "$(dirname "$0")"

# 检查应用是否在运行
# 使用 grep 查找进程，并排除 grep 本身
PID=$(ps -ef | grep "$APP_NAME" | grep -v grep | awk '{print $2}')

if [ -n "$PID" ]; then
    echo "正在停止 $APP_NAME (PID: $PID)..."
    kill $PID
    
    # 等待进程结束（最多等待 10 秒）
    for i in {1..10}; do
        if ps -p $PID > /dev/null; then
            sleep 1
        else
            break
        fi
    done
    
    # 如果还在运行，强制杀死
    if ps -p $PID > /dev/null; then
        echo "进程未响应，正在强制停止..."
        kill -9 $PID
    fi
    echo "$APP_NAME 已停止。"
else
    echo "$APP_NAME 未运行。"
fi

echo "正在启动 $APP_NAME..."

# 后台启动应用
# > startup.log 2>&1: 将标准输出和错误输出重定向到 startup.log，防止 nohup.out 膨胀，同时也便于排查启动报错
nohup ./$APP_NAME > startup.log 2>&1 &

# 等待几秒检查启动状态
sleep 2

# 获取新进程 ID
NEW_PID=$(ps -ef | grep "$APP_NAME" | grep -v grep | awk '{print $2}')

if [ -n "$NEW_PID" ]; then
    echo "$APP_NAME 启动成功! (PID: $NEW_PID)"
    echo "启动日志: startup.log"
else
    echo "$APP_NAME 启动失败!"
    echo "请查看 startup.log 获取错误信息。"
    tail -n 10 startup.log
fi
