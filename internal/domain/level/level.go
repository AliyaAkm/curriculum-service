package level

import "github.com/google/uuid"

type Level struct {
	ID   uuid.UUID `gorm:"column:id;primary_key"`
	Name string    `gorm:"column:name"`
	Code string    `gorm:"column:code"`
}

func (Level) TableName() string {
	return "course_levels"
}
