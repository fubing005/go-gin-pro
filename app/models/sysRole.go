package models

import (
	"shalabing-gin/global"

	"gorm.io/gorm"
)

//角色模型

type Role struct {
	ID
	Title       string `json:"title" gorm:"size:50;not null;unique;comment:角色名称"`
	Description string `json:"description" gorm:"size:255;comment:角色描述"`
	Status      int    `json:"status" gorm:"type:tinyint;default:1;comment:状态(1:正常 2:禁用)"`
	CreateBy    uint   `json:"create_by" gorm:"default:0;comment:创建者"`
	CreateUser  string `json:"create_user" gorm:"-"`
	UpdateBy    uint   `json:"update_by" gorm:"default:0;comment:更新者"`
	UpdateUser  string `json:"update_user" gorm:"-"`
	Timestamps
	SoftDeletes
}

func (r *Role) AfterFind(tx *gorm.DB) (err error) {
	var createUsername string
	global.App.DB.Model(&Manager{}).Where("id = ?", r.CreateBy).Select("username").Scan(&createUsername)
	r.CreateUser = createUsername
	var updateUsername string
	global.App.DB.Model(&Manager{}).Where("id = ?", r.UpdateBy).Select("username").Scan(&updateUsername)
	r.UpdateUser = updateUsername
	return
}

func (Role) TableName() string {
	return "sys_role"
}
