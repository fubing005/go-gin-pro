package models

import (
	"shalabing-gin/global"

	"gorm.io/gorm"
)

type Post struct {
	ID
	PostName   string `json:"post_name" gorm:"size:50;not null;unique;comment:岗位名称"`
	PostCode   string `json:"post_code" gorm:"size:50;not null;unique;comment:岗位编码"`
	Sort       int    `json:"sort" gorm:"type:int;default:0;comment:显示顺序"`
	Status     int    `json:"status" gorm:"type:tinyint;default:1;comment:状态(1:正常 2:禁用)"`
	Remark     string `json:"remark" gorm:"size:500;comment:备注"`
	CreateBy   uint   `json:"create_by" gorm:"type:int;comment:创建者ID"`
	CreateUser string `json:"create_user" gorm:"-"`
	UpdateBy   uint   `json:"update_by" gorm:"type:int;comment:更新者ID"`
	UpdateUser string `json:"update_user" gorm:"-"`
	Timestamps
	SoftDeletes
}

func (p *Post) AfterFind(tx *gorm.DB) (err error) {
	var createUsername string
	global.App.DB.Model(&Manager{}).Where("id = ?", p.CreateBy).Select("username").Scan(&createUsername)
	p.CreateUser = createUsername
	var updateUsername string
	global.App.DB.Model(&Manager{}).Where("id = ?", p.UpdateBy).Select("username").Scan(&updateUsername)
	p.UpdateUser = updateUsername
	return
}

func (Post) TableName() string {
	return "sys_post"
}
