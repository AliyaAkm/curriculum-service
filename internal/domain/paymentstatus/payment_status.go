package paymentstatus

import "github.com/google/uuid"

type PaymentStatus struct {
	ID   uuid.UUID `gorm:"column:id"`
	Name string    `gorm:"column:name"`
	Code string    `gorm:"column:code"`
}

func (PaymentStatus) TableName() string {
	return "payment_statuses"
}
