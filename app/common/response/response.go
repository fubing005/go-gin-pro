package response

import (
	"net/http"
	"os"
	"shalabing-gin/core/errors"
	"shalabing-gin/core/trans"
	"shalabing-gin/global"
	"time"

	"github.com/gin-gonic/gin"
)

// 响应结构体
type Response struct {
	ErrorCode int         `json:"error_code"` // 自定义错误码
	Message   string      `json:"message"`    // 信息
	Data      interface{} `json:"data"`       // 数据
}

// Success 响应成功 ErrorCode 为 0 表示成功
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		0,
		trans.Trans("common.成功"),
		data,
	})
}

// Fail 响应失败 ErrorCode 不为 0 表示失败
func Fail(c *gin.Context, errorCode int, msg string) {
	switch msg {
	case "业务逻辑错误":
		msg = trans.Trans("common.业务逻辑错误")
	case "请求参数错误":
		msg = trans.Trans("common.请求参数错误")
	case "登录授权失效":
		msg = trans.Trans("common.登录授权失效")
	}

	c.JSON(http.StatusOK, Response{
		errorCode,
		msg,
		nil,
	})
}

// FailByError 失败响应 返回自定义错误的错误码、错误信息
func FailByError(c *gin.Context, error errors.CustomError) {
	Fail(c, error.ErrorCode, error.ErrorMsg)
}

// ValidateFail 请求参数验证失败
func ValidateFail(c *gin.Context, msg string) {
	Fail(c, errors.Errors.ValidateError.ErrorCode, msg)
}

// BusinessFail 业务逻辑失败
func BusinessFail(c *gin.Context, msg string) {
	Fail(c, errors.Errors.BusinessError.ErrorCode, msg)
}

// TokenFail Token 失效
func TokenFail(c *gin.Context) {
	FailByError(c, errors.Errors.TokenError)
}

func ServerError(c *gin.Context, err interface{}) {
	msg := "Internal Server Error"
	// 非生产环境显示具体错误信息
	if global.App.Config.App.Env != "prod" && os.Getenv(gin.EnvGinMode) != gin.ReleaseMode {
		// if _, ok := err.(error); ok {
		// 	msg = err.(error).Error()
		// }
		if _, ok := err.(string); ok {
			msg = err.(string)
		}
	}
	c.JSON(http.StatusInternalServerError, Response{
		http.StatusInternalServerError,
		msg,
		nil,
	})
	c.Abort()
}

func TooManyRequests(ctx *gin.Context, msg string, data interface{}, extraData ...interface{}) {
	ctx.JSON(http.StatusTooManyRequests, gin.H{
		"error_code": http.StatusTooManyRequests,
		"message":    msg,
		"data":       data,
		"ext":        extraData,
		"time":       time.Now().Unix(),
	})
}
