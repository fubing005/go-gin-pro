# 运行环境
FROM docker.1ms.run/library/alpine:latest

# 安装必要的依赖（如果需要）
RUN apk add --no-cache libc6-compat

#RUN mkdir -p /app && chmod 755 /app

# 设置工作目录
WORKDIR /app

# 复制编译后的二进制文件
COPY ./dist/* /app

# 为二进制文件添加执行权限
RUN chmod +x /app/server

# 开放端口
EXPOSE 8080

# 启动应用
CMD ["./server"]
