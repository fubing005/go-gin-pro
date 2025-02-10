package routes

import (
	"shalabing-gin/app/controllers/admin"
	"shalabing-gin/app/middleware"
	"shalabing-gin/app/services"

	"github.com/gin-gonic/gin"
)

// SetApiGroupRoutes 定义 api 分组路由
func SetAdminGroupRoutes(router *gin.RouterGroup) {
	//请求日志
	// router.Use(middleware.RequestLogger("admin"))

	router.GET("/captcha", admin.LoginController{}.Captcha)
	// 用户登录
	router.POST("/auth/login", admin.LoginController{}.Login)
	// middleware.InitAdminAuthMiddleware
	authRouter := router.Group("").Use(middleware.JWTAuth(services.AdminGuardName), middleware.InitAdminAuthMiddleware)
	{
		authRouter.POST("/media/image_upload", admin.MediaController{}.ImageUpload)

		//后台首页数据
		authRouter.GET("/main/index", admin.MainController{}.Index)

		//管理员信息
		authRouter.GET("/manager/manager_info", admin.ManagerController{}.ManagerInfo)
		authRouter.GET("/manager/logout", admin.ManagerController{}.LogOut)

		//管理员路由
		authRouter.GET("/manager/common", admin.ManagerController{}.Common)
		authRouter.GET("/manager/index", admin.ManagerController{}.Index)
		authRouter.POST("/manager/do_add", admin.ManagerController{}.DoAdd)
		authRouter.GET("/manager/edit", admin.ManagerController{}.Edit)
		authRouter.PUT("/manager/do_edit", admin.ManagerController{}.DoEdit)
		authRouter.DELETE("/manager/delete", admin.ManagerController{}.Delete)

		//角色路由
		authRouter.GET("/role/common", admin.RoleController{}.Common)
		authRouter.GET("/role/index", admin.RoleController{}.Index)
		authRouter.POST("/role/do_add", admin.RoleController{}.DoAdd)
		authRouter.GET("/role/edit", admin.RoleController{}.Edit)
		authRouter.PUT("/role/do_edit", admin.RoleController{}.DoEdit)
		authRouter.DELETE("/role/delete", admin.RoleController{}.Delete)
		authRouter.GET("/role/permission_auth", admin.RoleController{}.PermissionAuth)
		authRouter.POST("/role/permission_do_auth", admin.RoleController{}.PermissionDoAuth)
		authRouter.GET("/role/dept_auth", admin.RoleController{}.DeptAuth)
		authRouter.POST("/role/dept_do_auth", admin.RoleController{}.DeptDoAuth)

		//权限路由
		authRouter.GET("/permission/common", admin.PermissionController{}.Common)
		authRouter.GET("/permission/index", admin.PermissionController{}.Index)
		authRouter.POST("/permission/do_add", admin.PermissionController{}.DoAdd)
		authRouter.GET("/permission/edit", admin.PermissionController{}.Edit)
		authRouter.PUT("/permission/do_edit", admin.PermissionController{}.DoEdit)
		authRouter.DELETE("/permission/delete", admin.PermissionController{}.Delete)

		//岗位管理
		authRouter.GET("/post/common", admin.PostController{}.Common)
		authRouter.GET("/post/index", admin.PostController{}.Index)
		authRouter.POST("/post/do_add", admin.PostController{}.DoAdd)
		authRouter.GET("/post/edit", admin.PostController{}.Edit)
		authRouter.PUT("/post/do_edit", admin.PostController{}.DoEdit)
		authRouter.DELETE("/post/delete", admin.PostController{}.Delete)

		//部门管理
		authRouter.GET("/dept/common", admin.DeptController{}.Common)
		authRouter.GET("/dept/index", admin.DeptController{}.Index)
		authRouter.POST("/dept/do_add", admin.DeptController{}.DoAdd)
		authRouter.GET("/dept/edit", admin.DeptController{}.Edit)
		authRouter.PUT("/dept/do_edit", admin.DeptController{}.DoEdit)
		authRouter.DELETE("/dept/delete", admin.DeptController{}.Delete)
	}
}
