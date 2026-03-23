package tag

import "github.com/google/uuid"

type Tag struct {
	ID   uuid.UUID `gorm:"column:id"`
	Name string    `gorm:"column:name"`
	Code string    `gorm:"column:code"`
}

func (Tag) TableName() string {
	return "course_tags"
}
