package coursepoint

import (
	"github.com/google/uuid"
	"time"
)

type UserCoursePoints struct {
	ID        uuid.UUID `gorm:"column:id;primary_key"`
	LessonID  uuid.UUID `gorm:"column:lesson_id"`
	UserID    uuid.UUID `gorm:"column:user_id"`
	XP        int64     `gorm:"column:xp"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}
type Leaderboard struct {
	Place  int       `gorm:"-"`
	UserID uuid.UUID `gorm:"column:user_id"`
	XP     int64     `gorm:"column:xp"`
}

func (UserCoursePoints) TableName() string {
	return "user_course_points"
}
