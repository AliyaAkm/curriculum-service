package durationcategory

import "github.com/google/uuid"

type DurationCategory struct {
	ID   uuid.UUID `gorm:"column:id;primary_key"`
	Name string    `gorm:"column:name"`
	Code string    `gorm:"column:code"`
}

func (DurationCategory) TableName() string {
	return "course_duration_categories"
}
