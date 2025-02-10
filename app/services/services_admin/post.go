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

type postService struct{}

var PostService = new(postService)

var postCount int64

func (postService *postService) Common() (data map[string]interface{}) {
	data = make(map[string]interface{}, 2)
	status := []Status{{Code: 1, Name: "启用"}, {Code: 2, Name: "禁用"}}
	data["status"] = status

	// managerList := []models.ManagerResponse{}
	// global.App.DB.Model(&models.Manager{}).Where("status = ?", 1).Select("id,username,nickname").Scan(&managerList)
	// data["manager"] = managerList
	return
}

// 获取所有岗位
func (postService *postService) GetPosts(params request.PageQuery) (data map[string]interface{}) {
	data = make(map[string]interface{}, 2)
	postList := []models.Post{}
	global.App.DB.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&postList)
	data["list"] = postList

	global.App.DB.Model(&models.Post{}).Count(&postCount)
	data["count"] = postCount

	return
}

// 添加岗位-执行添加
func (postService *postService) DoAdd(params request_admin.PostAdd, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		//实例化post
		post := models.Post{
			PostName: params.PostName,
			PostCode: params.PostCode,
			Sort:     params.Sort,
			Status:   params.Status,
			Remark:   params.Remark,
			CreateBy: uint(id),
		}

		err := global.App.DB.Create(&post).Error
		if err != nil {
			err = errors.New(trans.Trans("admin.post.添加岗位失败"))
			return err
		}

		// 创建岗位操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		err = CommonService.CreateAdminLog(uint(id), "岗位模块", "创建岗位", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), utils.StructToJsonString(post), ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 编辑岗位-获取岗位
func (postService *postService) Edit(params request_admin.PostEditDelete) (post *models.Post) {
	global.App.DB.Where("id = ?", params.ID).Find(&post)
	return
}

// 编辑岗位-执行编辑
func (postService *postService) DoEdit(params request_admin.PostEdit, c *gin.Context) (err error) {
	global.App.Log.Info("params", zapcore.Field{Type: zapcore.StringType, Key: "params", String: utils.StructToJsonString(params)})
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		//实例化post
		post := models.Post{}
		global.App.DB.Where("id = ?", params.ID).Find(&post)
		post.PostName = params.PostName
		post.PostCode = params.PostCode
		post.Sort = params.Sort
		post.Status = params.Status
		post.Remark = params.Remark
		post.UpdateBy = uint(id)

		err := global.App.DB.Save(&post).Error
		if err != nil {
			err = errors.New("admin.post.编辑岗位失败")
			return err
		}

		// 创建岗位操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		err = CommonService.CreateAdminLog(uint(id), "岗位模块", "编辑岗位", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), utils.StructToJsonString(post), ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}

// 删除岗位-执行删除
func (postService *postService) Delete(params request_admin.PostEditDelete, c *gin.Context) (err error) {
	return global.App.DB.Transaction(func(tx *gorm.DB) error {
		err := global.App.DB.Delete(&models.Post{}, params.ID).Error
		if err != nil {
			err = errors.New("admin.post.删除岗位失败")
			return err
		}

		// 创建岗位操作日志
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		value, _ := c.Keys["id"].(string)
		id, _ := strconv.ParseUint(value, 10, 64)
		err = CommonService.CreateAdminLog(uint(id), "岗位模块", "删除岗位", c.Request.Method, c.Request.URL.Path, utils.StructToJsonString(params), "", ip, userAgent, 1, 0)
		if err != nil {
			return err
		}

		return nil
	})
}
