package services_admin

import (
	"errors"
	"shalabing-gin/app/common/request/request_admin"
	"shalabing-gin/app/models"
	"shalabing-gin/core/trans"
	"shalabing-gin/global"
	"shalabing-gin/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

type permissionService struct{}

var PermissionService = new(permissionService)

func (permissionService *permissionService) Common() (data map[string]interface{}) {
	data = make(map[string]interface{}, 2)
	status := []Status{{Code: 1, Name: "启用"}, {Code: 2, Name: "禁用"}}
	data["status"] = status
	return
}

// 获取所有权限
func (permissionService *permissionService) GetPermissions() (permissionList []models.Permission) {
	global.App.DB.Where("module_id = ?", 0).Preload("PermissionItem.PermissionItem").Find(&permissionList)
	return
}

// 添加权限-执行添加
func (permissionService *permissionService) DoAdd(params request_admin.PermissionAdd, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)

		//实例化permission
		permission := models.Permission{
			ModuleName:  params.ModuleName,
			ActionName:  params.ActionName,
			Icon:        params.Icon,
			Type:        params.Type,
			Method:      params.Method,
			Url:         params.Url,
			ModuleId:    params.ModuleId,
			Sort:        params.Sort,
			Status:      params.Status,
			Description: params.Description,
			CreateBy:    uint(id),
		}

		err := global.App.DB.Create(&permission).Error
		if err != nil {
			err = errors.New(trans.Trans("admin.permission.添加权限失败"))
			return err
		}

		// 创建权限操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		err = CommonService.CreateAdminLog(uint(id), "权限模块", "创建权限", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), utils.StructToJsonString(permission), ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 编辑权限-获取权限
func (permissionService *permissionService) Edit(params request_admin.PermissionEditDelete) (permission *models.Permission) {
	global.App.DB.Where("id = ?", params.ID).Find(&permission)
	return
}

// 编辑权限-执行编辑
func (permissionService *permissionService) DoEdit(params request_admin.PermissionEdit, c *gin.Context) (err error) {
	global.App.Log.Info("params", zapcore.Field{Type: zapcore.StringType, Key: "params", String: utils.StructToJsonString(params)})
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		//实例化permission
		permission := models.Permission{}

		global.App.DB.Where("id = ?", params.ID).Find(&permission)
		permission.ModuleName = params.ModuleName
		permission.ActionName = params.ActionName
		permission.Icon = params.Icon
		permission.Type = params.Type
		permission.Method = params.Method
		permission.Url = params.Url
		permission.ModuleId = params.ModuleId
		permission.Sort = params.Sort
		permission.Status = params.Status
		permission.Description = params.Description
		permission.UpdateBy = uint(params.ID)

		err := global.App.DB.Save(&permission).Error
		if err != nil {
			err = errors.New("admin.permission.编辑权限失败")
			return err
		}

		// 创建权限操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		err = CommonService.CreateAdminLog(uint(id), "权限模块", "编辑权限", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), utils.StructToJsonString(permission), ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 删除权限-执行删除
func (permissionService *permissionService) Delete(params request_admin.PermissionEditDelete, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		err = global.App.DB.Delete(&models.Permission{}, params.ID).Error
		if err != nil {
			err = errors.New("admin.permission.删除权限失败")
			return err
		}

		err = global.App.DB.Where("permission_id = ?", params.ID).Delete(&models.RolePermission{}).Error
		if err != nil {
			err = errors.New("admin.role.删除关联角色失败")
			return err
		}

		// 创建权限操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		err = CommonService.CreateAdminLog(uint(id), "权限模块", "删除权限", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), "", ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}
