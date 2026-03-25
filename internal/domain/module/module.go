package module

import (
	"github.com/google/uuid"
	"time"
)

type Module struct {
	ID          uuid.UUID `gorm:"column:id"`
	CourseID    uuid.UUID `gorm:"course_id"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`
	Locale      string    `gorm:"locale"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Module) TableName() string {
	return "course_modules"
}
