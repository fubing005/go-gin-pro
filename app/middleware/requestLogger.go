package middleware

import (
	"shalabing-gin/global"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger 请求日志
func RequestLogger(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		defer func() {
			end := time.Now()
			latency := end.Sub(start).String()

			requestParams := make(map[string]interface{})

			for key, values := range c.Request.URL.Query() {
				if len(values) > 0 {
					requestParams[key] = values[0]
				}
			}

			err := c.Request.ParseForm()
			if err == nil {
				for key, values := range c.Request.PostForm {
					if len(values) > 0 {
						requestParams[key] = values[0]
					}
				}
			}

			clientIP := c.ClientIP()

			fields := []zap.Field{
				zap.String("clientIP", clientIP),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
				zap.Any("params", requestParams),
				zap.String("latency", latency),
			}
			global.App.Log.Info(version+":Request handled", fields...)
		}()
		c.Next()
	}
}
