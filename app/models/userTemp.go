package models

// 用户发送验证码相关
type UserTemp struct {
	ID
	Ip        string `json:"ip" gorm:"size:50;comment:操作IP"`
	Phone     string `json:"phone" gorm:"size:11;comment:手机号"`
	SendCount int    `json:"send_count" gorm:"type:int;comment:发送次数"`
	AddDay    string `json:"add_day" gorm:"size:10;comment:生成日期"`
	Sign      string `json:"sign" gorm:"type:varchar(255);comment:页面跳转标签"` //页面跳转标签
	Timestamps
	SoftDeletes
}

func (UserTemp) TableName() string {
	return "user_temp"
}
