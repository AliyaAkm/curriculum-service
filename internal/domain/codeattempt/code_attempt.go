package codeattempt

import (
	"time"

	"github.com/google/uuid"
)

const (
	RunTypeRun    = "run"
	RunTypeSubmit = "submit"
)

type Attempt struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	CourseID     *uuid.UUID
	LessonID     *uuid.UUID
	PracticeID   string
	RunType      string
	Language     string
	Passed       bool
	ErrorType    string
	ErrorMessage string
	Output       string
	DurationMS   int
	CodeHash     string
	XPReward     int
	XPAwarded    int
	CreatedAt    time.Time
}

type RunRequest struct {
	UserID     uuid.UUID
	CourseID   *uuid.UUID
	LessonID   *uuid.UUID
	PracticeID string
	RunType    string
	Language   string
	Code       string
}

type RunResult struct {
	AttemptID  uuid.UUID
	Output     string
	Error      string
	Passed     bool
	ErrorType  string
	DurationMS int
	XPAwarded  int
}
