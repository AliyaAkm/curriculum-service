package codeattempt

import "github.com/google/uuid"

type RunRequest struct {
	CourseID uuid.UUID `json:"course_id"`
	LessonID uuid.UUID `json:"lesson_id"`
	RunType  string    `json:"run_type"`
	Language string    `json:"language"`
	Code     string    `json:"code"`
}

type RunResponse struct {
	AttemptID  uuid.UUID `json:"attempt_id"`
	Output     string    `json:"output"`
	Error      string    `json:"error"`
	Passed     bool      `json:"passed"`
	ErrorType  string    `json:"error_type"`
	DurationMS int       `json:"duration_ms"`
	XPAwarded  int       `json:"xp_awarded"`
}
