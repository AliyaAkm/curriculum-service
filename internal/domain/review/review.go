package review

import (
	"curriculum-service/internal/domain/course"
	"github.com/google/uuid"
	"time"
)

type CourseReview struct {
	ID        uuid.UUID `gorm:"column:id;primary_key"`
	CourseID  uuid.UUID `gorm:"column:course_id"`
	UserID    uuid.UUID `gorm:"column:user_id"`
	User      course.User
	Rating    int       `gorm:"column:rating"`
	Comment   string    `gorm:"column:comment"`
	ViewCount int       `gorm:"column:view_count"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (CourseReview) TableName() string {
	return "course_reviews"
}
