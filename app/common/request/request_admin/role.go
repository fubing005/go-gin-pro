package request_admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/core/trans"
)

type RoleAdd struct {
	Title       string `form:"title" json:"title" binding:"required,min=3,max=20"`
	Description string `form:"description" json:"description" binding:"max=100"`
	Status      int    `form:"status" json:"status" binding:"required,oneof=1 2"`
}

func (roleAdd RoleAdd) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"title.required":  trans.Trans("admin.role.角色名称不能为空"),
		"status.required": trans.Trans("common.状态不能为空"),
		"status.oneof":    trans.Trans("common.状态不正确"),
	}
}

type RoleEditDelete struct {
	ID uint `form:"id" json:"id" binding:"required,exist_role"`
}

func (roleEditDelete RoleEditDelete) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":   trans.Trans("admin.role.角色ID不能为空"),
		"id.exist_role": trans.Trans("admin.超管无需授权"),
	}
}

type RoleEdit struct {
	ID          uint   `form:"id" json:"id" binding:"required,exist_role"`
	Title       string `form:"title" json:"title" binding:"required,min=3,max=20"`
	Description string `form:"description" json:"description" binding:"max=100"`
	Status      int    `form:"status" json:"status" binding:"required,oneof=1 2"`
}

func (roleEdit RoleEdit) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":     trans.Trans("admin.role.角色ID不能为空"),
		"id.exist_role":   trans.Trans("admin.role.角色不存在"),
		"title.required":  trans.Trans("admin.role.角色名称不能为空"),
		"status.required": trans.Trans("common.状态不能为空"),
		"status.oneof":    trans.Trans("common.状态不正确"),
	}
}

type RolePermissionAuth struct {
	RoleId         uint   `form:"role_id" json:"role_id" binding:"required,exist_role"`
	PermissionNode []uint `form:"permission_node[]" json:"permission_node[]" binding:"required,permission_slice"`
}

func (rolePermissionAuth RolePermissionAuth) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"role_id.required":                   trans.Trans("admin.role.角色ID不能为空"),
		"role_id.exist_role":                 trans.Trans("admin.role.角色不存在"),
		"permission_node[].required":         trans.Trans("admin.permission.权限ID不能为空"),
		"permission_node[].permission_slice": trans.Trans("admin.permission.权限ID格式不正确或不存在"),
	}
}

type RoleDeptAuth struct {
	RoleId   uint   `form:"role_id" json:"role_id" binding:"required,exist_role"`
	DeptNode []uint `form:"dept_node[]" json:"dept_node[]" binding:"required,dept_slice"`
}

func (roleDeptAuth RoleDeptAuth) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"role_id.required":       trans.Trans("admin.role.角色ID不能为空"),
		"role_id.exist_role":     trans.Trans("admin.role.角色不存在"),
		"dept_node[].required":   trans.Trans("admin.dept.部门ID不能为空"),
		"dept_node[].dept_slice": trans.Trans("admin.dept.部门ID格式不正确或不存在"),
	}
}
