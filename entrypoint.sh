#!/bin/sh
# 创建符号链接
ln -sf /app/config/config.json /app/config.json

# 运行主程序
exec ./main