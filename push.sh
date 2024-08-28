#!/bin/bash

# 服务器地址
SERVER="119.28.81.104"

# PEM 文件路径
PEM_FILE="./weshow_test.pem"

# 本地文件路径列表，可以根据实际情况添加更多文件路径
LOCAL_FILES=(
    "im"
    "config.json"
    "start.sh"
)

REMOTE_DIR="/www/goserver/chat"
REMOTE_SCRIPT="start.sh"

# 上传文件
for FILE in "${LOCAL_FILES[@]}"; do
    scp -i "$PEM_FILE" "$FILE" "root@$SERVER:$REMOTE_DIR"
done

# 连接到服务器并执行脚本
ssh -i "$PEM_FILE" root@"$SERVER" "cd $REMOTE_DIR &&./$REMOTE_SCRIPT"