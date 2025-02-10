package services_admin

import (
	"shalabing-gin/app/models"
	"shalabing-gin/global"
)

type commonService struct{}

var CommonService = new(commonService)

type Status struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

// CreateAdminLog 创建管理员操作日志
func (commonService *commonService) CreateAdminLog(adminID uint, module, action, method, url, params, result, ip, userAgent string, status int, duration int64) error {
	log := &models.ManagerLog{
		AdminID:   adminID,
		Module:    module,
		Action:    action,
		Method:    method,
		URL:       url,
		Params:    params,
		Result:    result,
		IP:        ip,
		UserAgent: userAgent,
		Status:    status,
		Duration:  duration,
	}
	return global.App.DB.Create(log).Error
}

// CreateAdminLoginLog 创建管理员登录日志
func (commonService *commonService) CreateAdminLoginLog(adminID uint, username, ip, location, device, browser, os string, status int, message string) error {
	log := &models.ManagerLoginLog{
		AdminID:  adminID,
		Username: username,
		IP:       ip,
		Location: location,
		Device:   device,
		Browser:  browser,
		OS:       os,
		Status:   status,
		Message:  message,
	}
	return global.App.DB.Create(log).Error
}
