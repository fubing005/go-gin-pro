package request_api

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/core/trans"
)

type Register struct {
	Nickname     string `form:"nickname" json:"nickname" binding:"required"`
	Mobile       string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password     string `form:"password" json:"password" binding:"required"`
	CaptchaId    string `form:"captcha_id" json:"captcha_id" binding:"required"`
	CaptchaValue string `form:"captcha_value" json:"captcha_value" binding:"required"`
}

// 自定义错误信息
func (register Register) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"captcha_id.required":    trans.Trans("API.验证码ID不能为空"),
		"captcha_value.required": trans.Trans("API.验证码不能为空"),
		"nickname.required":      trans.Trans("common.昵称不能为空"),
		"mobile.required":        trans.Trans("common.手机号码不能为空"),
		"mobile.mobile":          trans.Trans("common.手机号码格式不正确"),
		"password.required":      trans.Trans("common.密码不能为空"),
	}
}

type Login struct {
	Mobile       string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Password     string `form:"password" json:"password" binding:"required"`
	CaptchaId    string `form:"captcha_id" json:"captcha_id" binding:"required"`
	CaptchaValue string `form:"captcha_value" json:"captcha_value" binding:"required"`
}

func (login Login) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"captcha_id.required":    trans.Trans("common.验证码ID不能为空"),
		"captcha_value.required": trans.Trans("common.验证码不能为空"),
		"mobile.required":        trans.Trans("common.手机号码不能为空"),
		"mobile.mobile":          trans.Trans("common.手机号码格式不正确"),
		"password.required":      trans.Trans("common.密码不能为空"),
	}
}
