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
package api

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_api"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services"
	"shalabing-gin/app/services/services_api"
	"shalabing-gin/app/services/services_common"

	"github.com/gin-gonic/gin"
)

type LoginController struct{}

// captcha 获取验证码
// @Tags 获取验证码
// @Summary 获取验证码
// @Description 获取验证码接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/captcha [get]
func (con LoginController) Captcha(c *gin.Context) {
	id, b64s, _, err := services_common.MediaService.MakeCaptcha(50, 120, 4)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	outPut := map[string]interface{}{
		"captchaId":    id,
		"captchaImage": b64s,
	}

	response.Success(c, outPut)
}

// Register 用户注册
// @Tags 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Accept application/json
// @Produce application/json
// @Param data body request_api.Register true "用户注册"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/auth/register [post]
func (con LoginController) Register(c *gin.Context) {
	var form request_api.Register
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if user, err := services_api.UserService.Register(form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		response.Success(c, user)
	}
}

// Login 用户登录
// @Tags 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Accept application/json
// @Produce application/json
// @Param data body request_api.Login true "用户登录"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/auth/login [post]
func (con LoginController) Login(c *gin.Context) {
	var form request_api.Login
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if user, err := services_api.UserService.Login(form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		tokenData, err, _ := services.JwtService.CreateToken(services.ApiGuardName, user)
		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}
		response.Success(c, tokenData)
	}
}
