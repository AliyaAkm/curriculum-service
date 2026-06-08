package practice

import (
	"time"

	"github.com/google/uuid"
)

type TaskRequest struct {
	LessonID       uuid.UUID `json:"lesson_id"`
	Position       int       `json:"position"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Language       string    `json:"language"`
	StarterCode    string    `json:"starter_code"`
	ExpectedOutput string    `json:"expected_output"`
	XPReward       int       `json:"xp_reward"`
	CheckType      string    `json:"check_type"`
}

type TaskUpdateRequest struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	Language       *string `json:"language"`
	StarterCode    *string `json:"starter_code"`
	ExpectedOutput *string `json:"expected_output"`
	XPReward       *int    `json:"xp_reward"`
	CheckType      *string `json:"check_type"`
}

type TaskResponse struct {
	ID             uuid.UUID `json:"id"`
	LessonID       uuid.UUID `json:"lesson_id"`
	ModuleID       uuid.UUID `json:"module_id"`
	CourseID       uuid.UUID `json:"course_id"`
	Position       int       `json:"position"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Language       string    `json:"language"`
	StarterCode    string    `json:"starter_code"`
	ExpectedOutput string    `json:"expected_output,omitempty"`
	XPReward       int       `json:"xp_reward"`
	CheckType      string    `json:"check_type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
