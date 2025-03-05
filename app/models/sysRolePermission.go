package models

type RolePermission struct {
	RoleID       uint `json:"role_id" gorm:"column:role_id;not null;index:role_id;comment:角色ID"`
	PermissionID uint `json:"permission_id" gorm:"column:permission_id;not null;index:permission_id;comment:权限ID"`
}

func (RolePermission) TableName() string {
	return "sys_role_permission"
}
