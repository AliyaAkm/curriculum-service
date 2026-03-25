package locale

import "github.com/google/uuid"

type Locale struct {
	ID   uuid.UUID `gorm:"column:id"`
	Name string    `gorm:"column:name"`
	Code string    `gorm:"column:code"`
}

func (Locale) TableName() string {
	return "course_locales"
}
