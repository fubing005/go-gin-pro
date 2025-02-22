package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/middleware"
	"shalabing-gin/core/trans"
	"shalabing-gin/global"
	"shalabing-gin/routes"
	"syscall"
	"time"

	docs "shalabing-gin/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter() *gin.Engine {
	// router := gin.Default()

	if global.App.Config.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New() // 不使用默认中间件
	router.Use(gin.Logger(), middleware.CustomRecovery())

	// 全局限流
	if global.App.Config.App.Env == "prod" {
		router.Use(middleware.RateLimiter(global.App.Config.App.RequestLimit, time.Second*60))
	}

	// 跨域处理，项目语言切换
	router.Use(middleware.Cors(), middleware.LanguageMiddleware())

	// 找不到路由
	router.NoRoute(func(c *gin.Context) {
		response.ServerError(c, trans.Trans("common.请求地址不存在"))
	})
	// 找不到方法
	router.NoMethod(func(c *gin.Context) {
		response.ServerError(c, trans.Trans("common.请求方法不存在"))
	})

	// 前端项目静态资源
	router.StaticFile("/", "./static/dist/index.html")
	router.Static("/assets", "./static/dist/assets")
	router.StaticFile("/favicon.ico", "./static/dist/favicon.ico")
	// 其他静态资源
	router.Static("/public", "./static")
	router.Static("/storage", "./storage/app/public")

	//注册swagger文档
	registerSwagger(router)

	// 注册 api 分组路由
	apiGroup := router.Group("/api")
	routes.SetApiGroupRoutes(apiGroup)

	// 注册 admin 分组路由
	adminGroup := router.Group("/admin")
	routes.SetAdminGroupRoutes(adminGroup)

	// 注册websocket路由
	// router.GET("/ws", websocket.ChatController{}.ChatMessage)
	wsGroup := router.Group("/ws")
	routes.SetWebsocketGroupRoutes(wsGroup)

	return router
}

func registerSwagger(r gin.IRouter) {
	// API文档访问地址: http://host/swagger/index.html
	// 注解定义可参考 https://github.com/swaggo/swag#declarative-comments-format
	// 样例 https://github.com/swaggo/swag/blob/master/example/basic/api/api.go
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Title = "【管理后台/用户端】接口文档"
	docs.SwaggerInfo.Description = "实现前后端API服务"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// RunServer 启动服务器
func RunServer() {
	r := setupRouter()

	srv := &http.Server{
		Addr:    ":" + global.App.Config.App.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 启动消费者
	// go services_api.RabbitmqService.StartConsumer("queue_order", "exchange", "routing_key") //消费rabbitmq队列
	// go services_api.StartConsumer()//消费kafka队列
	// 消费订单状态数据，实现订单状态的变更,以及超时订单自动取消，
	// go services_api.ConsumeMessages()           // 启动消息队列消费者,用于处理订单状态
	// go services_api.ConsumeDeadLetterMessages() // 启动死信队列消费者,用于处理超时订单

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	closeService()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func closeService() {
	// 关闭数据库
	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()

	// 关闭Redis
	defer func() {
		if global.App.Redis != nil {
			global.App.Redis.Close()
		}
	}()

	// 关闭MongoDB
	defer func() {
		if global.App.MongoDB != nil {
			CloseMongoDB()
		}
	}()

	//关闭clickhouse
	defer func() {
		if global.App.Clickhouse != nil {
			global.App.Clickhouse.Close()
		}
	}()

	//关闭rabbitmq
	defer func() {
		if global.App.Rabbitmq != nil {
			global.App.Rabbitmq.GetConnection().Close()
		}
		if global.App.Kafka != nil {
			global.App.Kafka.GetProducer().Close()
		}
	}()
}
