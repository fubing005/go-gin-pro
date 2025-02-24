package models

import (
	"strconv"
)

type UserTimestamp struct {
}

type User struct {
	ID
	Mobile       string  `json:"mobile" gorm:"not null;unique;comment:用户手机号"`
	Password     string  `json:"-" gorm:"not null;default:'';comment:用户密码"`
	Username     string  `json:"username" gorm:"size:50;not null;unique;comment:用户名"`
	Nickname     string  `json:"nickname" gorm:"size:50;comment:昵称"`
	Avatar       string  `json:"avatar" gorm:"size:255;comment:头像"`
	Gender       int     `json:"gender" gorm:"type:tinyint;default:0;comment:性别(0:未知 1:男 2:女)"`
	Email        string  `json:"email" gorm:"size:100;comment:邮箱"`
	Status       int     `json:"status" gorm:"type:tinyint;default:1;comment:状态(1:正常 2:禁用)"`
	OrderCount   int64   `json:"order_count" gorm:"type:bigint;default:0;comment:订单数"`
	TotalAmount  float64 `json:"total_amount" gorm:"type:decimal(10,2);default:0;comment:总金额"`
	ViewCount    int64   `json:"view_count" gorm:"type:bigint;default:0;comment:浏览量"`
	FavorCount   int64   `json:"favor_count" gorm:"type:bigint;default:0;comment:收藏量"`
	CommentCount int64   `json:"comment_count" gorm:"type:bigint;default:0;comment:评论量"`
	Balance      float64 `json:"balance" gorm:"type:decimal(10,2);default:0;comment:余额"`
	RealName     string  `json:"real_name" gorm:"size:50;comment:真实姓名"`
	IDCard       string  `json:"id_card" gorm:"size:50;unique;comment:身份证号"`
	ThirdPartyID string  `json:"third_party_id" gorm:"size:100;unique;comment:第三方ID"`
	Provider     string  `json:"provider" gorm:"size:50;comment:第三方平台"`
	Birthday     MyTime  `json:"birthday" gorm:"default:NULL;comment:生日"`
	LastLogin    MyTime  `json:"last_login" gorm:"default:NULL;comment:最后登录时间"`
	Timestamps
	SoftDeletes
}

func (user User) GetUid() string {
	return strconv.Itoa(int(user.ID.ID))
}

func (User) TableName() string {
	return "users"
}
