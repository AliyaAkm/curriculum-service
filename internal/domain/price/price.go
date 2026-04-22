package price

import "github.com/google/uuid"

type CoursePrice struct {
	ID       uuid.UUID `gorm:"column:id"`
	CourseID uuid.UUID `gorm:"column:course_id"`
	Amount   int       `gorm:"column:amount"`
	Currency string    `gorm:"column:currency"`
}
