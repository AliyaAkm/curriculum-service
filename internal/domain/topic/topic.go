package topic

import "github.com/google/uuid"

type Topic struct {
	ID   uuid.UUID `gorm:"column:id"`
	Name string    `gorm:"column:name"`
	Code string    `gorm:"column:code"`
}

func (Topic) TableName() string {
	return "course_topics"
}
