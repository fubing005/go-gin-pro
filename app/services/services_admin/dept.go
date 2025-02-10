package services_admin

import (
	"errors"
	"shalabing-gin/app/common/request"
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

type deptService struct{}

var DeptService = new(deptService)

var deptCount int64

func (deptService *deptService) Common() (data map[string]interface{}) {
	data = make(map[string]interface{}, 2)
	status := []Status{{Code: 1, Name: "启用"}, {Code: 2, Name: "停用"}}
	data["status"] = status
	return
}

// 获取所有部门
func (deptService *deptService) GetDepts(form request.PageQuery) (data map[string]interface{}) {
	data = make(map[string]interface{}, 2)
	deptList := []models.Dept{}
	global.App.DB.Where("parent_id = ?", 0).Offset((form.Page - 1) * form.PageSize).Limit(form.PageSize).Preload("DeptItem").Find(&deptList)
	data["list"] = deptList

	global.App.DB.Model(&models.Dept{}).Where("parent_id = ?", 0).Count(&deptCount)
	data["count"] = deptCount

	return
}

// 添加部门-执行添加
func (deptService *deptService) DoAdd(params request_admin.DeptAdd, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		//实例化dept
		dept := models.Dept{
			ParentId: params.ParentId,
			DeptName: params.DeptName,
			Sort:     params.Sort,
			Leader:   params.Leader,
			Phone:    params.Phone,
			Email:    params.Email,
			Status:   params.Status,
			CreateBy: uint(id),
		}

		err := global.App.DB.Create(&dept).Error
		if err != nil {
			err = errors.New(trans.Trans("admin.添加部门失败"))
			return err
		}

		// 创建部门操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		err = CommonService.CreateAdminLog(uint(id), "部门模块", "创建部门", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), utils.StructToJsonString(dept), ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 编辑部门-获取部门
func (deptService *deptService) Edit(params request_admin.DeptEditDelete) (dept *models.Dept) {
	global.App.DB.Where("id = ?", params.ID).Find(&dept)
	return
}

// 编辑部门-执行编辑
func (deptService *deptService) DoEdit(params request_admin.DeptEdit, c *gin.Context) (err error) {
	global.App.Log.Info("params", zapcore.Field{Type: zapcore.StringType, Key: "params", String: utils.StructToJsonString(params)})
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		//实例化dept
		dept := models.Dept{}
		global.App.DB.Where("id = ?", params.ID).Find(&dept)
		dept.ParentId = params.ParentId
		dept.DeptName = params.DeptName
		dept.Sort = params.Sort
		dept.Leader = params.Leader
		dept.Phone = params.Phone
		dept.Email = params.Email
		dept.Status = params.Status
		err := global.App.DB.Save(&dept).Error
		if err != nil {
			err = errors.New("admin.编辑部门失败")
			return err
		}

		// 创建部门操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		err = CommonService.CreateAdminLog(uint(id), "部门模块", "编辑部门", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), utils.StructToJsonString(dept), ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 删除部门-执行删除
func (deptService *deptService) Delete(params request_admin.DeptEditDelete, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		err = global.App.DB.Delete(&models.Dept{}, params.ID).Error
		if err != nil {
			err = errors.New("admin.删除部门失败")
			return err
		}

		err := global.App.DB.Where("dept_id = ?", params.ID).Delete(&models.RoleDept{}).Error
		if err != nil {
			err = errors.New("admin.role.删除关联角色失败")
			return err
		}

		// 创建部门操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		err = CommonService.CreateAdminLog(uint(id), "部门模块", "删除部门", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), "", ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}
