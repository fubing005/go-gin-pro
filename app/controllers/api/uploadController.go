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
	"shalabing-gin/app/services/services_api"

	"shalabing-gin/app/common/response"

	"github.com/gin-gonic/gin"
)

type MediaController struct{}

// ImageUpload 用户上传文件
// @Summary 用户上传文件
// @Description 用户上传文件接口
// @Tags 文件上传
// @Accept multipart/form-data
// @Produce json
// @Security     Bearer
// @Param business formData string true "业务名称"
// @Param image formData file true "文件"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/media/image_upload [post]
func (con MediaController) ImageUpload(c *gin.Context) {
	var form request_api.ImageUpload
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	outPut, err := services_api.MediaService.SaveImage(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, outPut)
}
