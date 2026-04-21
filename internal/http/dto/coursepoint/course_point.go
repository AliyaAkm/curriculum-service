package coursepoint

import (
	"github.com/google/uuid"
	"time"
)

type CreateUserCoursePointsRequest struct {
	LessonID uuid.UUID `json:"lesson_id"`
	UserID   uuid.UUID `json:"user_id"`
	XP       int64     `json:"xp"`
}

type Leaderboard struct {
	Place  int       `json:"place"`
	UserID uuid.UUID `json:"user_id"`
	XP     int64     `json:"xp"`
}

type UpdateUserCoursePointsRequest struct {
	XP int64 `json:"xp"`
}

type UserCoursePointsResponse struct {
	ID        uuid.UUID `json:"id"`
	LessonID  uuid.UUID `json:"lesson_id"`
	UserID    uuid.UUID `json:"user_id"`
	XP        int64     `json:"xp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
