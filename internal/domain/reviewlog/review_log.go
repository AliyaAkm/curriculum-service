package reviewlog

import (
	"github.com/google/uuid"
	"time"
)

type CourseReviewLog struct {
	ID         uuid.UUID `gorm:"column:id;primary_key"`
	CourseID   uuid.UUID `gorm:"column:course_id"`
	AdminID    uuid.UUID `gorm:"column:admin_id"`
	IsApproved bool      `gorm:"column:is_approved"`
	Comment    string    `gorm:"column:comment"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (CourseReviewLog) TableName() string {
	return "course_review_logs"
}
