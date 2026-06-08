package practicereview

import (
	"time"

	"github.com/google/uuid"
)

type CreateSubmissionRequest struct {
	Code       string `json:"code"`
	Language   string `json:"language"`
	Output     string `json:"output"`
	Error      string `json:"error"`
	ErrorType  string `json:"error_type"`
	ExitCode   *int   `json:"exit_code"`
	DurationMS *int   `json:"duration_ms"`
}

type ReviewSubmissionRequest struct {
	Status  string `json:"status"`
	Comment string `json:"comment"`
}

type SubmissionResponse struct {
	ID                    uuid.UUID  `json:"id"`
	PracticeID            uuid.UUID  `json:"practice_id"`
	StudentID             uuid.UUID  `json:"student_id"`
	CourseID              uuid.UUID  `json:"course_id"`
	LessonID              uuid.UUID  `json:"lesson_id"`
	Status                string     `json:"status"`
	Code                  string     `json:"code"`
	Language              string     `json:"language"`
	Output                string     `json:"output"`
	Error                 string     `json:"error"`
	ErrorType             string     `json:"error_type"`
	ExitCode              *int       `json:"exit_code"`
	DurationMS            *int       `json:"duration_ms"`
	TeacherComment        string     `json:"teacher_comment"`
	ReviewedBy            *uuid.UUID `json:"reviewed_by"`
	ReviewedAt            *time.Time `json:"reviewed_at"`
	AttemptNumber         int        `json:"attempt_number"`
	PracticeTitle         string     `json:"practice_title"`
	StudentEmail          string     `json:"student_email"`
	CourseTitle           string     `json:"course_title"`
	LessonTitle           string     `json:"lesson_title"`
	ProgressStatus        string     `json:"progress_status"`
	ProgressStartedAt     *time.Time `json:"progress_started_at"`
	ProgressCompletedAt   *time.Time `json:"progress_completed_at"`
	ProgressLastAttemptAt *time.Time `json:"progress_last_attempt_at"`
	ProgressAttemptsCount int        `json:"progress_attempts_count"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}
