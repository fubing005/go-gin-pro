# 运行环境
FROM docker.1ms.run/library/alpine:latest

# 设置工作目录
WORKDIR /app

# 复制编译后的二进制文件
COPY ./dist /app

# 开放端口
EXPOSE 8080

# 启动应用
CMD ["./server"]
