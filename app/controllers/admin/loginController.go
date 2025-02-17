package admin

import (
	"fmt"
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services"
	"shalabing-gin/app/services/services_admin"
	"shalabing-gin/app/services/services_common"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
}

func (con LoginController) Captcha(c *gin.Context) {
	id, b64s, _, err := services_common.MediaService.MakeCaptcha(50, 120, 4)
	if err != nil {
		fmt.Println(err)
	}

	outPut := map[string]interface{}{
		"captchaId":    id,
		"captchaImage": b64s,
	}
	response.Success(c, outPut)
}

func (con LoginController) Login(c *gin.Context) {
	var form request_admin.Login
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if manager, err := services_admin.AdminService.Login(form, c); err != nil {
		response.BusinessFail(c, err.Error())
		return
	} else {
		tokenData, err, _ := services.JwtService.CreateToken(services.AdminGuardName, manager)
		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}
		response.Success(c, tokenData)
	}
}
