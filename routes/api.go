package routes

import (
	"shalabing-gin/app/controllers/api"
	"shalabing-gin/app/middleware"
	"shalabing-gin/app/services"

	"github.com/gin-gonic/gin"
)

// SetApiGroupRoutes 定义 api 分组路由
func SetApiGroupRoutes(router *gin.RouterGroup) {
	//请求日志
	router.Use(middleware.RequestLogger("api"))
	//获取验证码
	router.GET("/captcha", api.LoginController{}.Captcha)
	// 用户注册
	router.POST("/auth/register", api.LoginController{}.Register)
	// 用户登录
	router.POST("/auth/login", api.LoginController{}.Login)
	authRouter := router.Group("").Use(middleware.JWTAuth(services.ApiGuardName))
	{
		//上传图片
		authRouter.POST("/media/image_upload", api.MediaController{}.ImageUpload)
		// 获取用户信息
		authRouter.GET("/user/userinfo", api.UserController{}.UserInfo)
		// 用户退出
		authRouter.POST("/user/logout", api.UserController{}.Logout)
	}

	//rabbitmq
	//发布消息
	router.POST("/rabbitmq/publish_message", api.RabbitmqController{}.PublishMessage)
	router.POST("/rabbitmq/orders", api.RabbitmqController{}.CreateOrder)
	router.PUT("/rabbitmq/orders/status", api.RabbitmqController{}.UpdateOrderStatus)

	//kafka
	//发送消息
	router.POST("/kafka/send_message", api.KafkaController{}.SendMessage)

	//mongodb
	//获取文档
	router.GET("/mongodb/get_documents", api.MongoDBController{}.GetDocuments)
	router.POST("/mongodb/add_document", api.MongoDBController{}.AddDocument)
	router.PUT("/mongodb/update_document", api.MongoDBController{}.UpdateDocument)
	router.DELETE("/mongodb/delete_document", api.MongoDBController{}.DeleteDocument)
	router.GET("/mongodb/get_service_status", api.MongoDBController{}.GetServiceStatus)

	//clickhouse
	//获取文档
	router.POST("/clickhouse/create_table", api.ClickhouseController{}.CreateTable)
	router.POST("/clickhouse/insert", api.ClickhouseController{}.InsertHandler)
	router.GET("/clickhouse/query", api.ClickhouseController{}.QueryHandler)
	router.PUT("/clickhouse/update", api.ClickhouseController{}.UpdateHandler)
	router.DELETE("/clickhouse/delete", api.ClickhouseController{}.DeleteHandler)

	//elasticsearch
	router.GET("/elasticsearch/cat", api.ElasticsearchController{}.Cat)
	router.GET("/elasticsearch/get_index", api.ElasticsearchController{}.GetIndex)
	router.POST("/elasticsearch/create_index", api.ElasticsearchController{}.CreateIndex)
	router.DELETE("/elasticsearch/delete_index", api.ElasticsearchController{}.DeleteIndex)
	router.GET("/elasticsearch/get_documents", api.ElasticsearchController{}.GetDocuments)
	router.POST("/elasticsearch/add_document", api.ElasticsearchController{}.AddDocument)
	router.PUT("/elasticsearch/update_document", api.ElasticsearchController{}.UpdateDocument)
	router.DELETE("/elasticsearch/delete_document", api.ElasticsearchController{}.DeleteDocument)

	// 一些问题的解决方案

}
