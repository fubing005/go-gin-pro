/*
*

		_ooOoo_
	              o8888888o
	              88" . "88
	              (| -_- |)
	              O\  =  /O
	           ____/`---'\____
	         .'  \\|     |//  `.
	        /  \\|||  :  |||//  \
	       /  _||||| -:- |||||_  \
	       |   | \\\  -  /'| |   |
	       | \_|  `\`---'//  |_/ |
	       \  .-\__ `-. -'__/-.  /
	     ___`. .'  /--.--\  `. .'___
	  ."" '<  `.___\_<|>_/___.' _> \"".
	 | | :  `- \`. ;`. _/; .'/ /  .' ; |
	 \  \ `-.   \_\_`. _.'_/_/  -' _.' /

===================================

	佛祖加持 代码无BUG

===================================
*/
package admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services/services_admin"

	"github.com/gin-gonic/gin"
)

type RoleController struct{}

// 通用数据
func (con RoleController) Common(c *gin.Context) {
	commonData := services_admin.RoleService.Common()
	response.Success(c, commonData)
}

// 获取所有角色
func (con RoleController) Index(c *gin.Context) {
	var form request.PageQuery
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	roleList := services_admin.RoleService.Index(form, c)
	response.Success(c, roleList)
}

// 添加角色:提交
func (con RoleController) DoAdd(c *gin.Context) {
	var form request_admin.RoleAdd
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.RoleService.DoAdd(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 编辑角色:获取角色
func (con RoleController) Edit(c *gin.Context) {
	var form request_admin.RoleEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	role := services_admin.RoleService.Edit(form)
	response.Success(c, role)
}

// 编辑角色:提交
func (con RoleController) DoEdit(c *gin.Context) {
	var form request_admin.RoleEdit
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.RoleService.DoEdit(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 删除角色
func (con RoleController) Delete(c *gin.Context) {
	var form request_admin.RoleEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.RoleService.Delete(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 分配权限:获取权限
func (con RoleController) PermissionAuth(c *gin.Context) {
	var form request_admin.RoleEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	role, permissionList := services_admin.RoleService.PermissionAuth(form)

	data := map[string]interface{}{"role": role, "permission_list": permissionList}
	response.Success(c, data)
}

// 分配权限:提交
func (con RoleController) PermissionDoAuth(c *gin.Context) {
	var form request_admin.RolePermissionAuth
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.RoleService.PermissionDoAuth(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 分配部门:获取部门
func (con RoleController) DeptAuth(c *gin.Context) {
	var form request_admin.RoleEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	role, deptList := services_admin.RoleService.DeptAuth(form)

	data := map[string]interface{}{"role": role, "dept_list": deptList}
	response.Success(c, data)
}

// 分配部门:提交
func (con RoleController) DeptDoAuth(c *gin.Context) {
	var form request_admin.RoleDeptAuth
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.RoleService.DeptDoAuth(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}
