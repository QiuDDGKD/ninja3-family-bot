#!/bin/bash

# 定义进程号文件
PID_FILE="bot.pid"

# 如果 PID 文件存在，尝试杀掉进程及其子孙进程
if [ -f "$PID_FILE" ]; then
  PID=$(cat "$PID_FILE")
  if ps -p $PID > /dev/null 2>&1; then
    echo "Killing process $PID and its descendants..."
    pkill -TERM -P $PID  # 杀掉子孙进程
    kill -TERM $PID      # 杀掉主进程
    sleep 2              # 等待进程完全退出
  else
    echo "No process with PID $PID is running."
  fi
  rm -f "$PID_FILE"       # 删除 PID 文件
fi

# 启动程序并记录进程号
echo "Starting program..."
nohup go run main.go > bot.log 2>&1 &
NEW_PID=$!
echo $NEW_PID > "$PID_FILE"
echo "Program started with PID $NEW_PID."