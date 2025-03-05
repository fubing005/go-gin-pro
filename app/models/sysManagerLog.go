package models

// AdminLog 管理员操作日志
type ManagerLog struct {
	ID
	AdminID   uint   `json:"admin_id" gorm:"type:int;index;comment:管理员ID"`
	Module    string `json:"module" gorm:"size:50;comment:操作模块"`
	Action    string `json:"action" gorm:"size:50;comment:操作动作"`
	Method    string `json:"method" gorm:"size:20;comment:请求方法"`
	URL       string `json:"url" gorm:"size:200;comment:请求URL"`
	Params    string `json:"params" gorm:"type:text;comment:请求参数"`
	Result    string `json:"result" gorm:"type:text;comment:操作结果"`
	IP        string `json:"ip" gorm:"size:50;comment:操作IP"`
	UserAgent string `json:"user_agent" gorm:"size:200;comment:用户代理"`
	Status    int    `json:"status" gorm:"type:tinyint;comment:状态(1:成功 2:失败)"`
	Duration  int64  `json:"duration" gorm:"type:int;comment:执行时长(毫秒)"`
	Timestamps
	SoftDeletes
}

func (ManagerLog) TableName() string {
	return "sys_manager_log"
}
