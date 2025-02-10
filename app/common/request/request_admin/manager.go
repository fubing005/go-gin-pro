package request_admin

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/core/trans"
)

type Login struct {
	Username     string `form:"username" json:"username" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	CaptchaId    string `form:"captcha_id" json:"captcha_id" binding:"required"`
	CaptchaValue string `form:"captcha_value" json:"captcha_value" binding:"required"`
}

func (login Login) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"captcha_id.required":    trans.Trans("common.验证码ID不能为空"),
		"captcha_value.required": trans.Trans("common.验证码不能为空"),
		"username.required":      trans.Trans("common.账号不能为空"),
		"password.required":      trans.Trans("common.密码不能为空"),
	}
}

type ManagerAdd struct {
	Username   string `form:"username" json:"username" binding:"required,manager_username,exist_manager_username"`
	Password   string `form:"password" json:"password" binding:"required,manager_password"`
	RePassword string `form:"re_password" json:"re_password" binding:"required,eqfield=Password"`
	Nickname   string `form:"nickname" json:"nickname" binding:"omitempty,min=3,max=6"`
	Email      string `form:"email" json:"email" binding:"omitempty,email"`
	Mobile     string `form:"mobile" json:"mobile" binding:"omitempty,mobile"`
	RoleId     uint   `form:"role_id" json:"role_id" binding:"required,exist_role"`
	DeptId     uint   `form:"dept_id" json:"dept_id" binding:"required,exist_dept"`
	PostId     uint   `form:"post_id" json:"post_id" binding:"required,exist_post"`
	Status     int    `form:"status" json:"status" binding:"required,oneof=1 2"`
}

func (managerAdd ManagerAdd) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"username.required":               trans.Trans("common.账号不能为空"),
		"username.manager_username":       trans.Trans("admin.manager.管理员账号格式不正确"),
		"username.exist_manager_username": trans.Trans("admin.manager.管理员账号已存在"),
		"password.required":               trans.Trans("common.密码不能为空"),
		"password.manager_password":       trans.Trans("admin.manager.管理员密码格式不正确"),
		"re_password.required":            trans.Trans("common.确认密码不能为空"),
		"re_password.eqfield":             trans.Trans("common.两次密码不一致"),
		// "nickname.required":               trans.Trans("common.昵称不能为空"),
		// "email.required":                  trans.Trans("common.邮箱不能为空"),
		"email.email": trans.Trans("common.邮箱格式不正确"),
		// "mobile.required":                 trans.Trans("common.手机号码不能为空"),
		"mobile.mobile":      trans.Trans("common.手机号码格式不正确"),
		"role_id.required":   trans.Trans("admin.role.角色ID不能为空"),
		"role_id.exist_role": trans.Trans("admin.role.角色不存在"),
		"dept_id.required":   trans.Trans("admin.dept.部门ID不能为空"),
		"dept_id.exist_dept": trans.Trans("admin.dept.部门不存在"),
		"post_id.required":   trans.Trans("admin.post.岗位ID不能为空"),
		"post_id.exist_post": trans.Trans("admin.post.岗位不存在"),
		"status.required":    trans.Trans("common.状态不能为空"),
		"status.oneof":       trans.Trans("common.状态不正确"),
	}
}

type ManagerEditDelete struct {
	ID uint `form:"id" json:"id" binding:"required,exist_manager"`
}

func (managerEditDelete ManagerEditDelete) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":      trans.Trans("admin.manager.管理员ID不能为空"),
		"id.exist_manager": trans.Trans("admin.manager.管理员不存在"),
	}
}

type ManagerEdit struct {
	ID       uint   `form:"id" json:"id" binding:"required,exist_manager"`
	Username string `form:"username" json:"username" binding:"required,manager_username"`
	Password string `form:"password" json:"password" binding:"required,manager_password"`
	Nickname string `form:"nickname" json:"nickname" binding:"omitempty,min=3,max=6"`
	Email    string `form:"email" json:"email" binding:"omitempty,email"`
	Mobile   string `form:"mobile" json:"mobile" binding:"omitempty,mobile"`
	RoleId   uint   `form:"role_id" json:"role_id" binding:"required,exist_role"`
	DeptId   uint   `form:"dept_id" json:"dept_id" binding:"required,exist_dept"`
	PostId   uint   `form:"post_id" json:"post_id" binding:"required,exist_post"`
	Status   int    `form:"status" json:"status" binding:"required,oneof=1 2"`
}

func (managerEdit ManagerEdit) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":               trans.Trans("admin.manager.管理员ID不能为空"),
		"id.exist_manager":          trans.Trans("admin.manager.管理员不存在"),
		"username.required":         trans.Trans("common.账号不能为空"),
		"username.manager_username": trans.Trans("admin.manager.管理员账号格式不正确"),
		"password.required":         trans.Trans("common.密码不能为空"),
		"password.manager_password": trans.Trans("admin.manager.管理员密码格式不正确"),
		// "nickname.required":         trans.Trans("common.昵称不能为空"),
		// "email.required":            trans.Trans("common.邮箱不能为空"),
		"email.email": trans.Trans("common.邮箱格式不正确"),
		// "mobile.required":           trans.Trans("common.手机号码不能为空"),
		"mobile.mobile":      trans.Trans("common.手机号码格式不正确"),
		"role_id.required":   trans.Trans("admin.role.角色ID不能为空"),
		"role_id.exist_role": trans.Trans("admin.role.角色不存在"),
		"dept_id.required":   trans.Trans("admin.dept.部门ID不能为空"),
		"dept_id.exist_dept": trans.Trans("admin.dept.部门不存在"),
		"post_id.required":   trans.Trans("admin.post.岗位ID不能为空"),
		"post_id.exist_post": trans.Trans("admin.post.岗位不存在"),
		"status.required":    trans.Trans("common.状态不能为空"),
		"status.oneof":       trans.Trans("common.状态不正确"),
	}
}
