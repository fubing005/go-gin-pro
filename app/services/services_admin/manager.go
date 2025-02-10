package services_admin

import (
	"errors"
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/models"
	"shalabing-gin/app/services/services_common"
	"shalabing-gin/core/trans"
	"shalabing-gin/global"
	"shalabing-gin/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type adminService struct {
}

var AdminService = new(adminService)

var managerCount int64

func (adminService *adminService) Common() (data map[string]interface{}) {
	data = make(map[string]interface{}, 2)
	status := []Status{{Code: 1, Name: "启用"}, {Code: 2, Name: "禁用"}}
	data["status"] = status
	return
}

func (adminService *adminService) Login(params request_admin.Login, c *gin.Context) (manager *models.Manager, err error) {
	if flag := services_common.MediaService.VerifyCaptcha(params.CaptchaId, params.CaptchaValue); !flag {
		err = errors.New(trans.Trans("common.验证码不正确"))
		return
	}

	// 获取客户端信息
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	err = global.App.DB.Where("username = ?", params.Username).First(&manager).Error
	if err != nil || !utils.BcryptMakeCheck([]byte(params.Password), manager.Password) {
		err = errors.New(trans.Trans("common.用户名不存在或密码错误"))
		return
	}

	if manager.Status != 1 {
		_ = CommonService.CreateAdminLoginLog(0, params.Username, ip, "", "", "", "", 2, "账号已被禁用")
		err = errors.New(trans.Trans("common.账号已被禁用"))
		return
	}

	global.App.DB.Transaction(func(tx *gorm.DB) error {
		// 更新管理员最后登录时间
		manager.LastLogin = models.MyTime(time.Now())
		err = global.App.DB.Save(&manager).Error
		if err != nil {
			err = errors.New(trans.Trans("common.登录失败"))
			return err
		}

		device := utils.ParseUserAgent(userAgent)
		_ = CommonService.CreateAdminLoginLog(
			manager.ID.ID,
			manager.Username,
			ip,
			utils.GetLocationByIP(ip),
			device.Device,
			device.Browser,
			device.OS,
			1,
			"登录成功",
		)
		return nil
	})

	return
}

// GetUserInfo 获取管理员信息
func (adminService *adminService) GetAdminInfo(id string) (manager models.Manager, err error) {
	intId, _ := strconv.Atoi(id)
	err = global.App.DB.Preload("Role").First(&manager, intId).Error
	if err != nil {
		err = errors.New(trans.Trans("common.数据不存在"))
		return
	}
	return
}

func (adminService *adminService) Index(form request.PageQuery) (data map[string]interface{}) {
	data = make(map[string]interface{}, 2)
	managerList := []models.Manager{}
	global.App.DB.Offset((form.Page - 1) * form.PageSize).Limit(form.PageSize).Find(&managerList)
	data["list"] = managerList

	global.App.DB.Model(&models.Manager{}).Count(&managerCount)
	data["count"] = managerCount
	return
}

func (adminService *adminService) DoAdd(params request_admin.ManagerAdd, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		hash := utils.BcryptMake([]byte(params.Password))
		manager := models.Manager{
			Username: params.Username,
			Password: string(hash),
			Nickname: params.Nickname,
			Email:    params.Email,
			Mobile:   params.Mobile,
			RoleId:   params.RoleId,
			DeptId:   params.DeptId,
			PostId:   params.PostId,
			Status:   params.Status,
			CreateBy: uint(id),
		}
		err = global.App.DB.Create(&manager).Error
		if err != nil {
			err = errors.New(trans.Trans("admin.manager.添加管理员失败"))
			return err
		}

		// 创建管理员操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		err = CommonService.CreateAdminLog(uint(id), "管理员模块", "创建管理员", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), utils.StructToJsonString(manager), ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

func (adminService *adminService) Edit(params request_admin.ManagerEditDelete) (manager *models.Manager) {
	global.App.DB.Where("id = ?", params.ID).Find(&manager)
	return
}

func (adminService *adminService) DoEdit(params request_admin.ManagerEdit, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)

		manager := models.Manager{}
		// 检查用户是否存在
		global.App.DB.Where("id != ? AND username = ?", params.ID, params.Username).First(&manager)
		if manager.ID.ID > 0 {
			err = errors.New(trans.Trans("admin.manager.管理员账号已存在"))
			return err
		}

		global.App.DB.Find(&manager, params.ID)
		manager.Username = params.Username
		manager.Nickname = params.Nickname
		manager.Email = params.Email
		manager.Mobile = params.Mobile
		manager.RoleId = params.RoleId
		manager.DeptId = params.DeptId
		manager.PostId = params.PostId
		manager.Status = params.Status
		manager.UpdateBy = uint(id)

		hash := utils.BcryptMake([]byte(params.Password))
		manager.Password = string(hash)
		err = global.App.DB.Save(&manager).Error
		if err != nil {
			err = errors.New(trans.Trans("admin.manager.编辑管理员失败"))
			return err
		}

		// 创建管理员操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		err = CommonService.CreateAdminLog(uint(id), "管理员模块", "编辑管理员", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), "", ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

func (adminService *adminService) Delete(params request_admin.ManagerEditDelete, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		err = global.App.DB.Delete(&models.Manager{}, params.ID).Error
		if err != nil {
			err = errors.New(trans.Trans("admin.manager.删除管理员失败"))
			return err
		}

		// 创建管理员操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		err = CommonService.CreateAdminLog(uint(id), "管理员模块", "删除管理员", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), "", ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}
