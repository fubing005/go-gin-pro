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
	"shalabing-gin/app/services/services_admin"
	"shalabing-gin/utils"
	"strconv"

	"shalabing-gin/app/common/response"

	"github.com/gin-gonic/gin"
)

type MediaController struct{}

func (m MediaController) ImageUpload(c *gin.Context) {
	var form request_admin.ImageUpload
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	outPut, err := services_admin.MediaService.SaveImage(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	// 创建管理员操作日志
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	value, _ := c.Keys["id"].(string)
	id, _ := strconv.ParseUint(value, 10, 64)
	err = services_admin.CommonService.CreateAdminLog(uint(id), "管理员", "上传图片", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(form), utils.StructToJsonString(outPut), ip, userAgent, 1, 0)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, outPut)
}
