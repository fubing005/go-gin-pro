package middleware

//中间件: 作用: 在执行路由之前或者之后进行相关逻辑判断

import (
	"encoding/json"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/models"
	"shalabing-gin/core/trans"
	"shalabing-gin/global"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitAdminAuthMiddleware(c *gin.Context) {
	//权限判断: 没有登录的用户不能进入后台管理中心

	//1、获取Url访问的地址
	//当地址后面带参数时:,如: admin/captcha?t=0.8706946438889653,需要处理
	//strings.Split(c.Request.URL.String(), "?"): 把c.Request.URL.String()请求地址按照?分割成切片
	pathname := strings.Split(c.Request.URL.String(), "?")[0]

	//2、获取token上下文里面保存的用户信息
	userinfo := c.Keys["userinfo"]
	//3、判断Session中的用户信息是否存在，如果不存在跳转到登录页面（注意需要判断） 如果存在继续向下执行
	//session.Get获取返回的结果是一个空接口类型,所以需要进行类型断言: 判断userinfo是不是一个string
	userinfoStr, ok := userinfo.(string)
	if ok { // 说明是一个string
		var userinfoStruct models.Manager
		//把获取到的用户信息转换结构体
		err := json.Unmarshal([]byte(userinfoStr), &userinfoStruct)
		if err != nil || userinfoStruct.ID.ID == 0 { //表示未登录
			if pathname != "/admin/auth/login" && pathname != "/admin/captcha" {
				//跳转到登录页面
				response.TokenFail(c)
				c.Abort()
				return
			}
		} else { //表示用户登录成功
			//获取当前访问的URL对应的权限id,判断权限id是否在角色对应的权限中
			// strings.Replace 字符串替换
			urlPath := strings.Replace(pathname, "/admin/", "", 1)
			//排除权限判断:不是超级管理员并且不在相关权限内
			if userinfoStruct.IsSuper == 0 && !excludeAuthPath("/"+urlPath) {
				//判断用户权限:当前用户权限是否可以访问url地址
				//获取当前角色拥有的权限,并把权限id放在一个map对象中
				rolePermission := []models.RolePermission{}
				global.App.DB.Where("role_id = ?", userinfoStruct.RoleId).Find(&rolePermission)

				//判断当前用户所属部门，根据部门查询部门所属角色，将查询到的角色id放在rolePermission????
				roleIds := []uint{}
				global.App.DB.Model(&models.RoleDept{}).Where("dept_id = ?", userinfoStruct.DeptId).Select("role_id").Find(&roleIds)
				roleDeptPermission := []models.RolePermission{}
				global.App.DB.Where("role_id in (?)", roleIds).Find(&roleDeptPermission)
				rolePermission = append(rolePermission, roleDeptPermission...)

				rolePermissionMap := make(map[uint]uint)
				for _, v := range rolePermission {
					rolePermissionMap[v.PermissionID] = v.PermissionID
				}

				//实例化permission
				permission := models.Permission{}
				//查询权限id
				global.App.DB.Where("url = ? ", urlPath).Find(&permission)
				//判断权限id是否在角色对应的权限中
				if _, ok := rolePermissionMap[permission.ID.ID]; !ok {
					response.BusinessFail(c, trans.Trans("admin.permission.没有权限"))
					c.Abort() // 终止程序
					return
				}
			}
		}
	} else {
		//4、如果Session不存在，判断当前访问的URl是否是login doLogin captcha，如果不是跳转到登录页面，如果是不行任何操作
		//说明用户没有登录
		//需要排除到不需要做权限判断的路由
		if pathname != "/admin/auth/login" && pathname != "/admin/captcha" {
			//跳转到登录页面
			response.TokenFail(c)
			c.Abort()
			return
		}
	}
}

// 排除权限判断的方法
func excludeAuthPath(urlPath string) bool {
	//获取需要排除的地址
	excludeAuthPath := global.App.Config.Admin.ExcludeAuthPath
	//拆分字符串成为一个切片
	excludeAuthPathSlice := strings.Split(excludeAuthPath, ",")
	//判断传入的地址是否在排除地址内
	for _, v := range excludeAuthPathSlice {
		if v == urlPath {
			return true
		}
	}
	return false
}
