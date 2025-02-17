# 使用 golang:alpine 作为构建环境
FROM docker.1ms.run/library/golang:1.23 AS builder

# 安装依赖
#RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 先复制 go.mod 和 go.sum（加快缓存）
COPY go.mod go.sum ./

# 编译为 Linux 版可执行文件
RUN go mod tidy 

# 复制本地代码到容器
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gin_server main.go

# 运行环境（只保留编译后的二进制）
FROM docker.1ms.run/library/alpine:latest

WORKDIR /app
COPY --from=builder /app/gin_server /app/gin_server
# COPY --from=builder /app/configDocker.yaml ./config.yaml
# COPY --from=builder /app/static ./static/

RUN chmod +x /app/gin_server

EXPOSE 8080

# 运行应用
ENTRYPOINT  ["/app/gin_server"]
# CMD ["--config", "/app/config.yaml"]
