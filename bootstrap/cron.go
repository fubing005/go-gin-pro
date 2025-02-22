package bootstrap

import (
	"fmt"
	"log"
	"shalabing-gin/app/models"
	"shalabing-gin/app/services/services_api"
	"shalabing-gin/global"
	"time"

	"github.com/robfig/cron/v3"
)

func InitializeCron() {
	// global.App.Cron = cron.New(cron.WithSeconds())
	// go func() {
	// 	//每分钟执行一次
	// 	global.App.Cron.AddFunc("0 * * * * *", func() {
	// 		//使用 zapcore.Field 记录执行日志的时间
	// 		global.App.Log.Info("cron job run", zapcore.Field{Type: zapcore.StringType, Key: "cron", String: "every minute"})
	// 	})
	// 	global.App.Cron.Start()
	// 	defer global.App.Cron.Stop()
	// 	select {}
	// }()

	// 添加定时任务取消订单
	// go func() {
	// 	global.App.Cron = cron.New()
	// 	// 每分钟执行一次检查超时订单任务
	// 	global.App.Cron.AddFunc("@every 1m", func() {
	// 		var orders []models.Order
	// 		timeout := time.Now().Add(-15 * time.Minute) // 设置超时为 15 分钟
	// 		// 查找所有“待支付”且超时的订单
	// 		if err := global.App.DB.Where("status = ? AND created_at < ?", models.Pending, timeout).Find(&orders).Error; err != nil {
	// 			log.Println("查询超时订单失败:", err)
	// 			return
	// 		}
	// 		// 取消所有超时订单
	// 		for _, order := range orders {
	// 			order.Status = models.Canceled
	// 			if err := global.App.DB.Save(&order).Error; err != nil {
	// 				log.Println("取消订单失败:", order.ID)
	// 			} else {
	// 				log.Printf("订单 %d 已被取消，由于超时未支付\n", order.ID)
	// 			}
	// 		}
	// 	})
	// 	// 启动定时任务
	// 	global.App.Cron.Start()
	// }()

	//定时检查并取消超时订单
	go func() {
		global.App.Cron = cron.New()
		// 每分钟执行一次检查超时订单任务
		global.App.Cron.AddFunc("@every 1m", func() {
			var orders []models.Order
			timeout := time.Now().Add(-15 * time.Minute) // 设置超时为 15 分钟
			// 查找所有“待支付”且超时的订单
			if err := global.App.DB.Where("status = ? AND created_at < ?", models.Pending, timeout).Find(&orders).Error; err != nil {
				log.Println("查询超时订单失败:", err)
				return
			}
			// 将超时订单转发到死信队列
			for _, order := range orders {
				// 构造消息发送到死信队列
				msg := fmt.Sprintf(`{"order_id":%d, "status":"%s"}`, order.ID, models.Canceled)
				err := services_api.SendToDeadLetterQueue(msg)
				if err != nil {
					log.Println("无法发送消息到死信队列:", err)
				} else {
					// 更新订单状态为已取消
					order.Status = models.Canceled
					global.App.DB.Save(&order)
					log.Printf("订单 %d 超时，已发送到死信队列并更新为已取消状态\n", order.ID)
				}
			}
		})
		// 启动定时任务
		global.App.Cron.Start()
	}()
}
