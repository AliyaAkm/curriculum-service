package orderstatus

import "github.com/google/uuid"

type OrderStatus struct {
	ID   uuid.UUID `gorm:"column:id"`
	Name string    `gorm:"column:name"`
	Code string    `gorm:"column:code"`
}

func (OrderStatus) TableName() string {
	return "order_statuses"
}
