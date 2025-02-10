package request_admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/core/trans"
)

type PostAdd struct {
	PostName string `form:"post_name" json:"post_name" binding:"required"`
	PostCode string `form:"post_code" json:"post_code" binding:"required"`
	Sort     int    `form:"sort" json:"sort" binding:"min=1"`
	Status   int    `form:"status" json:"status" binding:"required,oneof=1 2"` //1启用 2禁用
	Remark   string `form:"remark" json:"remark" binding:"max=100"`
}

func (postAdd PostAdd) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"post_name.required": trans.Trans("admin.post.岗位名称不能为空"),
		"post_code.required": trans.Trans("admin.post.岗位编码不能为空"),
		"sort.min":           trans.Trans("admin.common.排序值不能小于1"),
		"status.required":    trans.Trans("common.状态不能为空"),
		"status.oneof":       trans.Trans("common.状态不正确"),
	}
}

type PostEditDelete struct {
	ID uint `form:"id" json:"id" binding:"required,exist_post"`
}

func (postEditDelete PostEditDelete) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":   trans.Trans("admin.post.岗位ID不能为空"),
		"id.exist_post": trans.Trans("admin.post.岗位不存在"),
	}
}

type PostEdit struct {
	ID       uint   `form:"id" json:"id" binding:"required,exist_post"`
	PostName string `form:"post_name" json:"post_name" binding:"required"`
	PostCode string `form:"post_code" json:"post_code" binding:"required"`
	Sort     int    `form:"sort" json:"sort" binding:"min=1"`
	Status   int    `form:"status" json:"status" binding:"required,oneof=1 2"` //1启用 2禁用
	Remark   string `form:"remark" json:"remark" binding:"max=100"`
}

func (postEdit PostEdit) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":        trans.Trans("admin.post.岗位ID不能为空"),
		"id.exist_post":      trans.Trans("admin.post.岗位不存在"),
		"post_name.required": trans.Trans("admin.post.岗位名称不能为空"),
		"post_code.required": trans.Trans("admin.post.岗位编码不能为空"),
		"sort.min":           trans.Trans("admin.common.排序值不能小于1"),
		"status.required":    trans.Trans("common.状态不能为空"),
		"status.oneof":       trans.Trans("common.状态不正确"),
	}
}
