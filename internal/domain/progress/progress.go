package progress

import (
	"time"

	"github.com/google/uuid"
)

type CourseProgress struct {
	CourseID           uuid.UUID
	UserID             uuid.UUID
	StartedAt          *time.Time
	LastActivityAt     *time.Time
	CompletedAt        *time.Time
	CurrentLessonID    *uuid.UUID
	TotalLessons       int
	CompletedLessons   int
	ProgressPercent    int
	CompletedLessonIDs []uuid.UUID
	PassedQuizIDs      []uuid.UUID
	Modules            []ModuleProgress
}

type ModuleProgress struct {
	ModuleID         uuid.UUID
	Position         int
	IsOpen           bool
	TotalLessons     int
	CompletedLessons int
}
