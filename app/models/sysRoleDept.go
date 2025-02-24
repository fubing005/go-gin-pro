package models

type RoleDept struct {
	RoleID uint `json:"role_id" gorm:"type:int;not null;comment:角色ID"`
	DeptID uint `json:"dept_id" gorm:"type:int;not null;comment:部门ID"`
}

func (RoleDept) TableName() string {
	return "sys_role_dept"
}
