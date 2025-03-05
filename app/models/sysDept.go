package models

import (
	"shalabing-gin/global"

	"gorm.io/gorm"
)

type Dept struct {
	ID
	ParentId   uint   `json:"parent_id" gorm:"type:int;comment:上级ID"`
	DeptName   string `json:"dept_name" gorm:"type:varchar(50);not null;comment:部门名称"`
	Sort       int    `json:"sort" gorm:"type:int;comment:排序"`
	Leader     string `json:"leader" gorm:"type:varchar(50);comment:负责人"`
	Phone      string `json:"phone" gorm:"type:varchar(11);comment:手机号"`
	Email      string `json:"email" gorm:"type:varchar(50);comment:邮箱"`
	Status     int    `json:"status" gorm:"type:tinyint;default:1;comment:状态(1:正常 2:停用)"`
	CreateBy   uint   `json:"create_by" gorm:"default:0;comment:创建者"`
	CreateUser string `json:"create_user" gorm:"-"`
	UpdateBy   uint   `json:"update_by" gorm:"default:0;comment:更新者"`
	UpdateUser string `json:"update_user" gorm:"-"`
	DeptItem   []Dept `json:"dept_item" gorm:"foreignKey:ParentId;references:ID"` // 表的自关联,获取该数据的子分类
	Checked    bool   `json:"checked" gorm:"-"`                                   // 用户是否有该权, 忽略本字段,给struct加一个自定义属性,和数据库没有关系
	Timestamps
	SoftDeletes
}

func (p *Dept) AfterFind(tx *gorm.DB) (err error) {
	var createUsername string
	global.App.DB.Model(&Manager{}).Where("id = ?", p.CreateBy).Select("username").Scan(&createUsername)
	p.CreateUser = createUsername
	var updateUsername string
	global.App.DB.Model(&Manager{}).Where("id = ?", p.UpdateBy).Select("username").Scan(&updateUsername)
	p.UpdateUser = updateUsername
	return
}

func (Dept) TableName() string {
	return "sys_dept"
}
