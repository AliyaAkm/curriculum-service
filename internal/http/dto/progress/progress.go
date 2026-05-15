package progress

import (
	"time"

	"github.com/google/uuid"
)

type CourseProgress struct {
	CourseID           uuid.UUID        `json:"course_id"`
	UserID             uuid.UUID        `json:"user_id"`
	StartedAt          *time.Time       `json:"started_at,omitempty"`
	LastActivityAt     *time.Time       `json:"last_activity_at,omitempty"`
	CompletedAt        *time.Time       `json:"completed_at,omitempty"`
	CurrentLessonID    *uuid.UUID       `json:"current_lesson_id,omitempty"`
	TotalLessons       int              `json:"total_lessons"`
	CompletedLessons   int              `json:"completed_lessons"`
	ProgressPercent    int              `json:"progress_percent"`
	CompletedLessonIDs []uuid.UUID      `json:"completed_lesson_ids"`
	PassedQuizIDs      []uuid.UUID      `json:"passed_quiz_ids"`
	Modules            []ModuleProgress `json:"modules"`
}

type ModuleProgress struct {
	ModuleID         uuid.UUID `json:"module_id"`
	Position         int       `json:"position"`
	IsOpen           bool      `json:"is_open"`
	TotalLessons     int       `json:"total_lessons"`
	CompletedLessons int       `json:"completed_lessons"`
}
