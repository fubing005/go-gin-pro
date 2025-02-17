package admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services"
	"shalabing-gin/app/services/services_admin"
	"shalabing-gin/core/trans"
	"shalabing-gin/utils"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ManagerController struct {
}

func (con ManagerController) ManagerInfo(c *gin.Context) {
	user, err := services.JwtService.GetUserInfo(services.AdminGuardName, c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, user)
}

// LogOut 用户退出
// @Tags 后台管理系统
// @Summary 用户退出
// @Description 用户退出接口
// @Accept application/json
// @Produce application/json
// @Security     Bearer
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/manager/logout [post]
func (con ManagerController) LogOut(c *gin.Context) {
	err := services.JwtService.JoinBlackList(c.Keys["token"].(*jwt.Token))
	if err != nil {
		response.BusinessFail(c, trans.Trans("common.登出失败"))
		return
	}

	// 获取客户端信息
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	device := utils.ParseUserAgent(userAgent)
	value, _ := c.Keys["id"].(string)
	id, _ := strconv.ParseUint(value, 10, 64)
	err = services_admin.CommonService.CreateAdminLoginLog(
		uint(id),
		c.Keys["username"].(string),
		ip,
		utils.GetLocationByIP(ip),
		device.Device,
		device.Browser,
		device.OS,
		1, // 1:成功
		"退出登录成功",
	)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (con ManagerController) Common(c *gin.Context) {
	commonData := services_admin.AdminService.Common()
	response.Success(c, commonData)
}

func (con ManagerController) Index(c *gin.Context) {
	var form request.PageQuery
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	managerList := services_admin.AdminService.Index(form)
	response.Success(c, managerList)
}

func (con ManagerController) DoAdd(c *gin.Context) {
	var form request_admin.ManagerAdd
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err := services_admin.AdminService.DoAdd(form, c); err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (con ManagerController) Edit(c *gin.Context) {
	var form request_admin.ManagerEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	manager := services_admin.AdminService.Edit(form)
	response.Success(c, manager)
}

func (con ManagerController) DoEdit(c *gin.Context) {
	var form request_admin.ManagerEdit
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err := services_admin.AdminService.DoEdit(form, c); err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (con ManagerController) Delete(c *gin.Context) {
	var form request_admin.ManagerEditDelete
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err := services_admin.AdminService.Delete(form, c); err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nil)
}
