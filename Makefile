# 项目名称
PROJECT_NAME := server
# 源文件名
SOURCE_FILE := main.go
# 二进制文件名称
BINARY_NAME := $(PROJECT_NAME)
# 静态文件目录
STATIC_DIR := static
# Docker 镜像名称
DOCKER_IMAGE_NAME := golang_gin_server
# Docker 镜像标签
DOCKER_IMAGE_TAG := latest
#源配置文件
SOURCE_CONFIG_FILE := configDocker.yaml
#目标配置文件
DEST_CONFIG_FILE := config.yaml
#打包后的路径
DIST_DIR := dist

# 默认目标
all: build

# 构建应用
build: $(SOURCE_FILE) $(STATIC_DIR)
	@echo "Starting to build the application..."
    # 设置环境变量以交叉编译到 Linux 环境
	set GOOS=linux 
	set GOARCH=amd64
	set CGO_ENABLED=0
	go build -o $(BINARY_NAME) -ldflags="-s -w" $(SOURCE_FILE)
    # 创建打包目录
	mkdir $(DIST_DIR)
    # 移动二进制文件到打包目录
	move $(BINARY_NAME) $(DIST_DIR)/
    # 复制静态文件到打包目录
	xcopy /E /I /Y $(STATIC_DIR) $(DIST_DIR)\$(STATIC_DIR)
	copy  $(SOURCE_CONFIG_FILE) $(DIST_DIR)\$(DEST_CONFIG_FILE)
	@echo "Application build completed."

# 构建并运行应用
run: build
	@echo "Starting the application..."
    # 在本地运行应用
	./$(DIST_DIR)/$(BINARY_NAME)
    # 在 Docker 中运行应用
    # docker run -p 8080:8080 --rm $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) 
	@echo "Application stopped."

# 清理生成的二进制文件
clean:
	@echo "Cleaning up..."
    # 删除打包目录
	rmdir /S /Q $(DIST_DIR)
	@echo "Cleanup completed."

.PHONY: all build run clean latest