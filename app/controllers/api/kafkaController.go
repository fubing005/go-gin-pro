package api

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_api"
	"shalabing-gin/app/common/response"
	"shalabing-gin/global"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type KafkaController struct{}

func (con KafkaController) SendMessage(c *gin.Context) {
	var form request_api.KafkaRequest
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	message := &sarama.ProducerMessage{
		Topic: global.App.Config.Queue.Kafka.Topic,
		Key:   sarama.StringEncoder(form.Key),
		Value: sarama.StringEncoder(form.Value),
	}

	partition, offset, err := global.App.Kafka.GetProducer().SendMessage(message)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message":   "Message sent successfully",
		"partition": partition,
		"offset":    offset,
	})
}
