package status

import "github.com/google/uuid"

type Status struct {
	ID   uuid.UUID `gorm:"column:id;primary_key"`
	Name string    `gorm:"column:name"`
	Code string    `gorm:"column:code"`
}

func (Status) TableName() string {
	return "course_statuses"
}
