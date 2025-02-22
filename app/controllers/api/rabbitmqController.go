package api

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_api"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/models"
	"shalabing-gin/global"

	"shalabing-gin/app/services/services_api"

	"github.com/gin-gonic/gin"
)

type RabbitmqController struct{}

// 测试rabbitmq发送消息
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

// 创建订单
func (con RabbitmqController) CreateOrder(c *gin.Context) {
	var form request_api.RabbitMQRequestOrder
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	form.Status = models.Pending // 初始状态

	order := models.Order{
		Amount: form.Amount,
		Status: form.Status,
	}
	if err := global.App.DB.Create(&order).Error; err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "订单创建失败"})
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	response.Success(c, order)
}

// 更新订单状态
func (con RabbitmqController) UpdateOrderStatus(c *gin.Context) {
	var form request_api.RabbitMQRequestOrder
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	var formID request_api.RabbitMQRequestOrderID
	if err := c.ShouldBindQuery(&formID); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(formID, err))
		return
	}

	var order models.Order
	if err := global.App.DB.First(&order, formID.ID).Error; err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err := order.UpdateStatus(form.NewStatus); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	// 将状态变更事件发送到 RabbitMQ
	go services_api.RabbitmqProducerService.PublishOrderUpdate(formID.ID, string(form.NewStatus))

	global.App.DB.Save(&order)
	response.Success(c, order)
}
