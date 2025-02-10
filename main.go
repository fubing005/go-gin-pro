package main

import (
	"shalabing-gin/bootstrap"
)

// @title       Gin Framework
// @version     1.0
// @description GIN框架API文档

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8080
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in                         header
// @name                      Authorization
// @description              Bearer token authentication
func main() {
	// 初始化服务
	bootstrap.InitializeService()

	// 启动服务器
	bootstrap.RunServer()
}
