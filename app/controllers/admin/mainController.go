package admin

import (
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services/services_admin"

	"github.com/gin-gonic/gin"
)

type MainController struct{}

func (con MainController) Index(c *gin.Context) {
	data, err := services_admin.MainService.Index(c)
	if err != nil {
		response.TokenFail(c)
		return
	}
	response.Success(c, data)
}
