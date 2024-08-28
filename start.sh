mv im showChat

chmod +x showChat

killall showChat

# 判断当前目录下是否存在 nohup.log 文件
if [ -f nohup.log ]; then
    # 将 nohup.log 的内容追加到 nohupbak.log
    cat nohup.log >> nohupbak.log
    # 删除 nohup.log 文件
    rm nohup.log
    echo "操作完成。"
else
    echo "当前目录下不存在 nohup.log 文件。"
fi

nohup ./showChat >> nohup.log 2>&1 &

echo "start complete"
