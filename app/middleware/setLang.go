package middleware

import (
	"shalabing-gin/global"
	"shalabing-gin/utils"

	"github.com/gin-gonic/gin"
)

// 中间件：从请求头中获取语言设置
func LanguageMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")

		flag := utils.ContainStr([]string{"zh-cn", "en"}, lang)
		if !flag {
			lang = "zh-cn"
		}

		global.App.Config.App.Lang = lang
	}
}
