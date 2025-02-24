package models

import (
	"errors"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"   // 待支付
	Paid      OrderStatus = "paid"      // 已支付
	Shipped   OrderStatus = "shipped"   // 已发货
	Completed OrderStatus = "completed" // 已完成
	Canceled  OrderStatus = "canceled"  // 已取消
)

type Order struct {
	ID
	Amount    float64     `json:"amount" gorm:"type:decimal(10,2);default:0;comment:总金额"`
	Status    OrderStatus `json:"status" gorm:"type:varchar(20);default:'pending';comment:订单状态"`
	UserID    uint        `json:"user_id" gorm:"type:int;not null;comment:用户ID"`
	ProductID uint        `json:"product_id" gorm:"type:int;not null;comment:产品ID"`
	Timestamps
}

func (o *Order) TableName() string {
	return "orders"
}

// 状态机：状态变更逻辑
func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	validTransitions := map[OrderStatus][]OrderStatus{
		Pending:   {Paid, Canceled},
		Paid:      {Shipped},
		Shipped:   {Completed},
		Completed: {},
		Canceled:  {},
	}

	// 检查状态是否可以变更
	validNextStates, exists := validTransitions[o.Status]
	if !exists {
		return errors.New("无效的状态转换")
	}

	for _, s := range validNextStates {
		if s == newStatus {
			o.Status = newStatus
			return nil
		}
	}

	return errors.New("非法的状态变更")
}
