package middleware

import (
	"shalabing-gin/app/common/response"
	"shalabing-gin/core/trans"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 限流中间件
func RateLimiter(maxRequests int, timeWindow time.Duration) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(timeWindow/time.Duration(maxRequests)), maxRequests)
	return func(c *gin.Context) {
		if limiter.Allow() == false {
			response.TooManyRequests(c, trans.Trans("common.请求频繁"), nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
