package models

// LoginLog 登录日志
type ManagerLoginLog struct {
	ID
	AdminID  uint   `json:"admin_id" gorm:"type:int;index;comment:管理员ID"`
	Username string `json:"username" gorm:"size:50;index;comment:用户名"`
	IP       string `json:"ip" gorm:"size:50;comment:登录IP"`
	Location string `json:"location" gorm:"size:100;comment:登录地点"`
	Device   string `json:"device" gorm:"size:50;comment:登录设备"`
	Browser  string `json:"browser" gorm:"size:50;comment:浏览器"`
	OS       string `json:"os" gorm:"size:50;comment:操作系统"`
	Status   int    `json:"status" gorm:"type:tinyint;comment:状态(1:成功 2:失败)"`
	Message  string `json:"message" gorm:"size:200;comment:消息"`
	Timestamps
	SoftDeletes
}

func (ManagerLoginLog) TableName() string {
	return "sys_manager_login_log"
}
