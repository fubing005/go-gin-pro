package request_admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/core/trans"
)

type DeptAdd struct {
	ParentId uint   `form:"parent_id" json:"parent_id" binding:"omitempty,exist_dept"`
	DeptName string `form:"dept_name" json:"dept_name" binding:"required"`
	Sort     int    `form:"sort" json:"sort" binding:"min=1"`
	Leader   string `form:"leader" json:"leader" binding:"required"`
	Phone    string `form:"phone" json:"phone" binding:"required,mobile"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Status   int    `form:"status" json:"status" binding:"required,oneof=1 2"`
}

func (deptAdd DeptAdd) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"parent_id.exist_dept": trans.Trans("admin.dept.上级部门不存在"),
		"dept_name.required":   trans.Trans("admin.dept.部门名称不能为空"),
		"sort.min":             trans.Trans("admin.common.排序值不能小于1"),
		"leader.required":      trans.Trans("admin.dept.部门负责人不能为空"),
		"phone.required":       trans.Trans("common.手机号码不能为空"),
		"phone.mobile":         trans.Trans("common.手机号码格式不正确"),
		"email.required":       trans.Trans("common.邮箱不能为空"),
		"email.email":          trans.Trans("common.邮箱格式错误"),
		"status.required":      trans.Trans("common.状态不能为空"),
		"status.oneof":         trans.Trans("common.状态不正确"),
	}
}

type DeptEditDelete struct {
	ID uint `form:"id" json:"id" binding:"required,exist_dept"`
}

func (deptEditDelete DeptEditDelete) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":   trans.Trans("admin.dept.部门ID不能为空"),
		"id.exist_dept": trans.Trans("admin.dept.部门不存在"),
	}
}

type DeptEdit struct {
	ID       uint   `form:"id" json:"id" binding:"required,exist_dept"`
	ParentId uint   `form:"parent_id" json:"parent_id" binding:"omitempty,exist_dept"`
	DeptName string `form:"dept_name" json:"dept_name" binding:"required"`
	Sort     int    `form:"sort" json:"sort" binding:"min=1"`
	Leader   string `form:"leader" json:"leader" binding:"required"`
	Phone    string `form:"phone" json:"phone" binding:"required,mobile"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Status   int    `form:"status" json:"status" binding:"required,oneof=1 2"`
}

func (deptEdit DeptEdit) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":          trans.Trans("admin.dept.部门ID不能为空"),
		"id.exist_dept":        trans.Trans("admin.dept.部门不存在"),
		"parent_id.exist_dept": trans.Trans("admin.dept.上级部门不存在"),
		"dept_name.required":   trans.Trans("admin.dept.部门名称不能为空"),
		"sort.min":             trans.Trans("admin.common.排序值不能小于1"),
		"leader.required":      trans.Trans("admin.dept.部门负责人不能为空"),
		"phone.required":       trans.Trans("common.手机号码不能为空"),
		"phone.mobile":         trans.Trans("common.手机号码格式不正确"),
		"email.required":       trans.Trans("common.邮箱不能为空"),
		"email.email":          trans.Trans("common.邮箱格式错误"),
		"status.required":      trans.Trans("common.状态不能为空"),
		"status.oneof":         trans.Trans("common.状态不正确"),
	}
}
