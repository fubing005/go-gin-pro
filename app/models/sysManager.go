package models

import (
	"shalabing-gin/app/common/response/response_admin"
	"shalabing-gin/global"
	"strconv"

	"gorm.io/gorm"
)

// 管理员表
type Manager struct { // 结构体首字母大写, 和数据库表名对应, 默认访问数据表users, 可以设置访问数据表的方法
	ID
	Username   string                      `json:"username" gorm:"size:50;not null;unique;comment:用户名"`
	Password   string                      `json:"-" gorm:"not null;default:'';comment:用户密码"`
	Nickname   string                      `json:"nickname" gorm:"size:50;comment:昵称"`
	Mobile     string                      `json:"mobile" gorm:"not null;comment:手机号"`
	Email      string                      `json:"email" gorm:"size:100;comment:邮箱"`
	Status     int                         `json:"status" gorm:"type:tinyint;default:1;comment:状态(1:正常 2:禁用)"`
	RoleId     uint                        `json:"role_id" gorm:"type:int;comment:角色ID"`
	DeptId     uint                        `json:"dept_id" gorm:"type:int;comment:部门ID"`
	PostId     uint                        `json:"post_id" gorm:"type:int;comment:岗位ID"`
	IsSuper    int                         `json:"is_super" gorm:"type:tinyint;default:0;comment:是否超级管理员(0:否 1:是)"`
	LastLogin  MyTime                      `json:"last_login" gorm:"default:NULL;comment:最后登录时间"`
	Role       response_admin.RoleResponse `json:"role" gorm:"-;foreignKey:RoleId;references:ID"` // 配置关联关系
	Dept       response_admin.DeptResponse `json:"dept" gorm:"-;foreignKey:DeptId;references:ID"` // 配置关联关系
	Post       response_admin.PostResponse `json:"post" gorm:"-;foreignKey:PostId;references:ID"` // 配置关联关系
	CreateBy   uint                        `json:"create_by" gorm:"default:0;comment:创建者"`
	CreateUser string                      `json:"create_user" gorm:"-"`
	UpdateBy   uint                        `json:"update_by" gorm:"default:0;comment:更新者"`
	UpdateUser string                      `json:"update_user" gorm:"-"`
	Timestamps
	SoftDeletes
}

func (p *Manager) AfterFind(tx *gorm.DB) (err error) {
	createUsername := ""
	global.App.DB.Model(&Manager{}).Where("id = ?", p.CreateBy).Select("username").Scan(&createUsername)
	p.CreateUser = createUsername

	updateUsername := ""
	global.App.DB.Model(&Manager{}).Where("id = ?", p.UpdateBy).Select("username").Scan(&updateUsername)
	p.UpdateUser = updateUsername

	role := response_admin.RoleResponse{}
	global.App.DB.Model(&Role{}).Where("id = ?", p.RoleId).Select("id,title").Scan(&role)
	p.Role = role

	dept := response_admin.DeptResponse{}
	global.App.DB.Model(&Dept{}).Where("id = ?", p.DeptId).Select("id,dept_name").Scan(&dept)
	p.Dept = dept

	post := response_admin.PostResponse{}
	global.App.DB.Model(&Post{}).Where("id = ?", p.PostId).Select("id,post_name,post_code").Scan(&post)
	p.Post = post
	return
}

func (Manager) TableName() string {
	return "sys_manager"
}

func (manager Manager) GetUid() string {
	return strconv.Itoa(int(manager.ID.ID))
}
