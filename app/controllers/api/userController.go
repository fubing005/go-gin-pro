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
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services"
	"shalabing-gin/core/trans"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UserController struct{}

// UserInfo 获取用户信息
// @Tags 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息接口
// @Accept application/json
// @Produce application/json
// @Security     Bearer
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/user/info [get]
func (con UserController) UserInfo(c *gin.Context) {
	user, err := services.JwtService.GetUserInfo(services.ApiGuardName, c.Keys["id"].(string))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, user)
}

// Logout 用户退出
// @Tags 用户退出
// @Summary 用户退出
// @Description 用户退出接口
// @Accept application/json
// @Produce application/json
// @Security     Bearer
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/user/logout [post]
func (con UserController) Logout(c *gin.Context) {
	err := services.JwtService.JoinBlackList(c.Keys["token"].(*jwt.Token))
	if err != nil {
		response.BusinessFail(c, trans.Trans("common.登出失败"))
		return
	}
	response.Success(c, nil)
}
