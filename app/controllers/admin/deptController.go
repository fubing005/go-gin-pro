package admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services/services_admin"

	"github.com/gin-gonic/gin"
)

type DeptController struct{}

// 部门通用数据
func (con DeptController) Common(c *gin.Context) {
	commonData := services_admin.DeptService.Common()
	response.Success(c, commonData)
}

// 获取所有部门
func (con DeptController) Index(c *gin.Context) {
	var form request.PageQuery
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	deptList := services_admin.DeptService.GetDepts(form)
	response.Success(c, deptList)
}

// 添加部门-执行添加
func (con DeptController) DoAdd(c *gin.Context) {
	var form request_admin.DeptAdd
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.DeptService.DoAdd(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 编辑部门-获取单个部门、顶级部门
func (con DeptController) Edit(c *gin.Context) {
	var form request_admin.DeptEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	dept := services_admin.DeptService.Edit(form)
	response.Success(c, dept)
}

// 编辑部门-执行编辑
func (con DeptController) DoEdit(c *gin.Context) {
	var form request_admin.DeptEdit
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_admin.DeptService.DoEdit(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// 删除部门
func (con DeptController) Delete(c *gin.Context) {
	var form request_admin.DeptEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	err := services_admin.DeptService.Delete(form, c)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}
