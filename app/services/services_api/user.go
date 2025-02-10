package services_api

import (
	"errors"
	"shalabing-gin/app/common/request/request_api"
	"shalabing-gin/app/models"
	"shalabing-gin/app/services/services_common"
	"shalabing-gin/core/trans"
	"shalabing-gin/global"
	"shalabing-gin/utils"
	"strconv"
	"time"
)

type userService struct {
}

var UserService = new(userService)

// Register 注册
func (userService *userService) Register(params request_api.Register) (user models.User, err error) {
	if flag := services_common.MediaService.VerifyCaptcha(params.CaptchaId, params.CaptchaValue); !flag {
		err = errors.New(trans.Trans("common.验证码不正确"))
		return
	}

	var result = global.App.DB.Where("mobile = ?", params.Mobile).Select("id").First(&models.User{})
	if result.RowsAffected != 0 {
		err = errors.New(trans.Trans("common.手机号已存在"))
		return
	}
	user = models.User{Nickname: params.Nickname, Mobile: params.Mobile, Password: utils.BcryptMake([]byte(params.Password))}
	err = global.App.DB.Create(&user).Error
	return
}

// Login 登录
func (userService *userService) Login(params request_api.Login) (user *models.User, err error) {
	if flag := services_common.MediaService.VerifyCaptcha(params.CaptchaId, params.CaptchaValue); !flag {
		err = errors.New(trans.Trans("common.验证码不正确"))
		return
	}
	err = global.App.DB.Where("mobile = ?", params.Mobile).First(&user).Error
	if err != nil || !utils.BcryptMakeCheck([]byte(params.Password), user.Password) {
		err = errors.New(trans.Trans("common.用户名不存在或密码错误"))
		return
	}
	user.LastLogin = models.MyTime(time.Now())
	err = global.App.DB.Save(&user).Error
	if err != nil {
		err = errors.New(trans.Trans("common.登录失败"))
		return
	}
	return
}

// GetUserInfo 获取用户信息
func (userService *userService) GetUserInfo(id string) (user models.User, err error) {
	intId, _ := strconv.Atoi(id)

	err = global.App.DB.First(&user, intId).Error
	if err != nil {
		err = errors.New(trans.Trans("common.数据不存在"))

	}
	return
}
