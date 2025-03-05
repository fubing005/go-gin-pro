package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 自定义时间类型
type MyTime time.Time

// 自增ID主键
type ID struct {
	ID uint `json:"id" gorm:"primaryKey"`
}

// 创建、更新时间
type Timestamps struct {
	CreatedAt MyTime `json:"created_at"`
	UpdatedAt MyTime `json:"-"`
}

// 软删除
type SoftDeletes struct {
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// 自定义 JSON 序列化方法
func (t MyTime) MarshalJSON() ([]byte, error) {
	formattedTime := time.Time(t).Format("2006-01-02 15:04:05")
	return []byte(`"` + formattedTime + `"`), nil
}

// 自定义 JSON 反序列化方法
func (t *MyTime) UnmarshalJSON(data []byte) error {
	parsedTime, err := time.Parse(`"2006-01-02 15:04:05"`, string(data))
	if err != nil {
		return err
	}
	*t = MyTime(parsedTime)
	return nil
}

func (t *Timestamps) BeforeCreate(tx *gorm.DB) (err error) {
	if t.CreatedAt == (MyTime{}) { // 检查是否为零值
		t.CreatedAt = MyTime(time.Now())
	}
	if t.UpdatedAt == (MyTime{}) { // 检查是否为零值
		t.UpdatedAt = MyTime(time.Now())
	}
	return nil
}

func (t *Timestamps) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = MyTime(time.Now())
	return nil
}

// 转换为数据库的值（实现 driver.Valuer 接口）
func (t MyTime) Value() (driver.Value, error) {
	return time.Time(t).Format("2006-01-02 15:04:05"), nil
}

// 从数据库扫描值（实现 sql.Scanner 接口）
func (t *MyTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time: // 数据库返回标准时间类型
		*t = MyTime(v)
	case string: // 数据库返回字符串类型
		parsedTime, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		*t = MyTime(parsedTime)
	case []uint8: // 数据库返回字节切片类型
		parsedTime, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		*t = MyTime(parsedTime)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}
