package admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services/services_admin"

	"github.com/gin-gonic/gin"
)

type PermissionController struct{}

// 通用数据
func (con PermissionController) Common(c *gin.Context) {
	commonData := services_admin.PermissionService.Common()
	response.Success(c, commonData)
}

// 获取所有权限
func (con PermissionController) Index(c *gin.Context) {
	permissionList := services_admin.PermissionService.GetPermissions()
	response.Success(c, permissionList)
}

// 添加权限-执行添加
func (con PermissionController) DoAdd(c *gin.Context) {
	var form request_admin.PermissionAdd
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.PermissionService.DoAdd(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 编辑权限-获取单个权限、顶级权限
func (con PermissionController) Edit(c *gin.Context) {
	var form request_admin.PermissionEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	permission := services_admin.PermissionService.Edit(form)
	response.Success(c, permission)
}

// 编辑权限-执行编辑
func (con PermissionController) DoEdit(c *gin.Context) {
	var form request_admin.PermissionEdit
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.PermissionService.DoEdit(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 删除权限
func (con PermissionController) Delete(c *gin.Context) {
	var form request_admin.PermissionEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	err := services_admin.PermissionService.Delete(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}
