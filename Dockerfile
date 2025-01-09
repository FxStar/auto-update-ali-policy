# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将 Go 模块文件复制到工作目录
COPY go.mod .
COPY go.sum .

# 下载依赖
RUN go mod tidy

# 将整个项目复制到工作目录
COPY . .

# 构建应用程序
RUN go build -o main .

# 使用一个更小的基础镜像来运行应用程序
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从 builder 阶段复制构建好的二进制文件
COPY --from=builder /app/main .

# 复制启动脚本
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# 设置容器启动时执行的命令
ENTRYPOINT ["/entrypoint.sh"]