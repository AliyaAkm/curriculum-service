package practice

import (
	"time"

	"github.com/google/uuid"
)

const (
	CheckTypeAuto   = "auto"
	CheckTypeManual = "manual"
)

type Task struct {
	ID             uuid.UUID
	LessonID       uuid.UUID
	ModuleID       uuid.UUID
	CourseID       uuid.UUID
	Position       int
	Title          string
	Description    string
	Language       string
	StarterCode    string
	ExpectedOutput string
	XPReward       int
	CheckType      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type TaskUpdate struct {
	Title          *string
	Description    *string
	Language       *string
	StarterCode    *string
	ExpectedOutput *string
	XPReward       *int
	CheckType      *string
}
