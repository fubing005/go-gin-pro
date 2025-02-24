package models

import (
	"shalabing-gin/global"

	"gorm.io/gorm"
)

//权限模型

type Permission struct {
	ID
	ModuleName     string       `json:"module_name" gorm:"size:50;comment:模块名称"`
	ActionName     string       `json:"action_name" gorm:"size:50;comment:操作名称"`
	Icon           string       `json:"icon" gorm:"size:50;comment:图标"`
	Type           int          `json:"type" gorm:"type:tinyint;comment:节点类型[1：模块,2：菜单，3：操作]"`
	Method         string       `json:"method" gorm:"size:50;comment:请求方式[GET/POST/PUT/DELETE]"`
	Url            string       `json:"url" gorm:"size:200;comment:路由跳转地址"`
	ModuleId       uint         `json:"module_id" gorm:"type:int;comment:模块ID[此module_id和当前模型的id关联,module_id= 0 表示模块]"`
	Sort           int          `json:"sort" gorm:"type:int;comment:排序"`
	Description    string       `json:"description" gorm:"size:255;comment:描述"`
	Status         int          `json:"status" gorm:"type:tinyint;default:1;comment:状态(1:显示 2:隐藏)"`
	CreateBy       uint         `json:"create_by" gorm:"default:0;comment:创建者"`
	CreateUser     string       `json:"create_user" gorm:"-"`
	UpdateBy       uint         `json:"update_by" gorm:"default:0;comment:更新者"`
	UpdateUser     string       `json:"update_user" gorm:"-"`
	PermissionItem []Permission `json:"permission_item" gorm:"foreignKey:ModuleId;references:ID"` // 表的自关联,获取该数据的子分类
	Checked        bool         `json:"checked" gorm:"-"`                                         // 用户是否有该权, 忽略本字段,给struct加一个自定义属性,和数据库没有关系
	Timestamps
	SoftDeletes
}

func (p *Permission) AfterFind(tx *gorm.DB) (err error) {
	var createUsername string
	global.App.DB.Model(&Manager{}).Where("id = ?", p.CreateBy).Select("username").Scan(&createUsername)
	p.CreateUser = createUsername
	var updateUsername string
	global.App.DB.Model(&Manager{}).Where("id = ?", p.UpdateBy).Select("username").Scan(&updateUsername)
	p.UpdateUser = updateUsername
	return
}

func (Permission) TableName() string {
	return "sys_permission"
}
