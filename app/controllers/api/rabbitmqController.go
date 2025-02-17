package api

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_api"
	"shalabing-gin/app/common/response"

	"shalabing-gin/app/services/services_api"

	"github.com/gin-gonic/gin"
)

type RabbitmqController struct{}

func (con RabbitmqController) PublishMessage(c *gin.Context) {
	var form request_api.RabbitMQRequest
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := services_api.RabbitmqService.PublishMessage(form.Exchange, form.RoutingKey, form.Message)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, nil)
}
