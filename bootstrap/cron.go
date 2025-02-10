package bootstrap

import (
	"shalabing-gin/global"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap/zapcore"
)

func InitializeCron() {
	global.App.Cron = cron.New(cron.WithSeconds())

	go func() {
		//每分钟执行一次
		global.App.Cron.AddFunc("0 * * * * *", func() {
			//使用 zapcore.Field 记录执行日志的时间
			global.App.Log.Info("cron job run", zapcore.Field{Type: zapcore.StringType, Key: "cron", String: "every minute"})
		})
		global.App.Cron.Start()
		defer global.App.Cron.Stop()
		select {}
	}()
}
