package api

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_api"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/models"
	"shalabing-gin/global"
	"time"

	"shalabing-gin/app/services/services_api"

	"github.com/gin-gonic/gin"
)

type RabbitmqController struct{}

var redisCtx = context.Background()

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
	var form request_api.RabbitMQRequestOrderCreate
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	// 检查幂等性必填字段
	if form.UserID == 0 || form.ProductID == 0 {
		response.ValidateFail(c, request.GetErrorMsg(form, errors.New("缺少 user_id 或 product_id")))
		return
	}
	// 生成 `request_id`
	requestID := generateRequestID(form.UserID, form.ProductID)

	// 检查 Redis 是否已有
	exists, _ := global.App.Redis.Exists(redisCtx, requestID).Result()
	if exists > 0 {
		response.ValidateFail(c, request.GetErrorMsg(form, errors.New("请勿重复提交订单")))
		return
	}
	// 记录 `request_id`，防止重复提交
	err := global.App.Redis.Set(redisCtx, requestID, "1", 1*time.Minute).Err()
	if err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	form.Status = models.Pending // 初始状态
	order := models.Order{
		Amount:    form.Amount,
		Status:    form.Status,
		UserID:    form.UserID,
		ProductID: form.ProductID,
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
	var form request_api.RabbitMQRequestOrderStatusUpdate
	if err := c.ShouldBindQuery(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	var order models.Order
	if err := global.App.DB.First(&order, form.ID).Error; err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if err := order.UpdateStatus(form.NewStatus); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	// 将状态变更事件发送到 RabbitMQ
	go services_api.RabbitmqProducerService.PublishOrderUpdate(form.ID, string(form.NewStatus))
	response.Success(c, order)
}

// 获取幂等性ID
func generateRequestID(userID, productID uint) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%d_%d", userID, productID)))
	return hex.EncodeToString(hash[:]) // 返回 MD5 作为 `request_id`
}
