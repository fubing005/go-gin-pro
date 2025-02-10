package request_admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/core/trans"
)

type PermissionAdd struct {
	ModuleName  string `form:"module_name" json:"module_name" binding:"required"`
	ActionName  string `form:"action_name" json:"action_name" binding:"omitempty"`
	Icon        string `form:"icon" json:"icon" binding:"omitempty"`
	Type        int    `form:"type" json:"type" binding:"required,oneof=1 2 3"`
	Method      string `form:"method" json:"method" binding:"omitempty,oneof=GET POST PUT PATCH DELETE"`
	Url         string `form:"url" json:"url"`
	ModuleId    uint   `form:"module_id" json:"module_id" binding:"omitempty,exist_permission"`
	Sort        int    `form:"sort" json:"sort" binding:"min=1"`
	Status      int    `form:"status" json:"status" binding:"required,oneof=1 2"` //1启用 2禁用
	Description string `form:"description" json:"description" binding:"max=100"`
}

func (permissionAdd PermissionAdd) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"module_name.required":       trans.Trans("admin.permission.模块名称不能为空"),
		"type.required":              trans.Trans("admin.permission.节点类型不能为空"),
		"type.oneof":                 trans.Trans("admin.permission.节点类型不正确"),
		"method.oneof":               trans.Trans("admin.permission.请求方法不正确"),
		"module_id.exist_permission": trans.Trans("admin.权限不存在"),
		"sort.min":                   trans.Trans("admin.common.排序值不能小于1"),
		"status.required":            trans.Trans("common.状态不能为空"),
		"status.oneof":               trans.Trans("common.状态不正确"),
	}
}

type PermissionEditDelete struct {
	ID uint `form:"id" json:"id" binding:"required,exist_permission"`
}

func (permissionEditDelete PermissionEditDelete) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":         trans.Trans("admin.permission.权限ID不能为空"),
		"id.exist_permission": trans.Trans("admin.权限不存在"),
	}
}

type PermissionEdit struct {
	ID          uint   `form:"id" json:"id" binding:"required,exist_permission"`
	ModuleName  string `form:"module_name" json:"module_name" binding:"required"`
	ActionName  string `form:"action_name" json:"action_name" binding:"omitempty"`
	Icon        string `form:"icon" json:"icon" binding:"omitempty"`
	Type        int    `form:"type" json:"type" binding:"required,oneof=1 2 3"`
	Method      string `form:"method" json:"method" binding:"omitempty,oneof=GET POST PUT PATCH DELETE"`
	Url         string `form:"url" json:"url"`
	ModuleId    uint   `form:"module_id" json:"module_id" binding:"omitempty,exist_permission"`
	Sort        int    `form:"sort" json:"sort" binding:"min=1"`
	Status      int    `form:"status" json:"status" binding:"required,oneof=1 2"` //1启用 2禁用
	Description string `form:"description" json:"description" binding:"max=100"`
}

func (permissionEdit PermissionEdit) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":                trans.Trans("admin.permission.权限ID不能为空"),
		"id.exist_permission":        trans.Trans("admin.权限不存在"),
		"module_name.required":       trans.Trans("admin.permission.模块名称不能为空"),
		"type.required":              trans.Trans("admin.permission.节点类型不能为空"),
		"type.oneof":                 trans.Trans("admin.permission.节点类型不正确"),
		"method.oneof":               trans.Trans("admin.permission.请求方法不正确"),
		"module_id.exist_permission": trans.Trans("admin.权限不存在"),
		"sort.min":                   trans.Trans("admin.common.排序值不能小于1"),
		"status.required":            trans.Trans("common.状态不能为空"),
		"status.oneof":               trans.Trans("common.状态不正确"),
	}
}
