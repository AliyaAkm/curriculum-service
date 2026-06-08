package practicereview

import (
	"time"

	"github.com/google/uuid"
)

const (
	SubmissionStatusSubmitted        = "submitted"
	SubmissionStatusInReview         = "in_review"
	SubmissionStatusChangesRequested = "changes_requested"
	SubmissionStatusApproved         = "approved"

	ProgressStatusInProgress       = "in_progress"
	ProgressStatusSubmitted        = "submitted"
	ProgressStatusChangesRequested = "changes_requested"
	ProgressStatusCompleted        = "completed"
)

type Submission struct {
	ID                    uuid.UUID
	PracticeID            uuid.UUID
	StudentID             uuid.UUID
	CourseID              uuid.UUID
	LessonID              uuid.UUID
	Status                string
	Code                  string
	Language              string
	Output                string
	Error                 string
	ErrorType             string
	ExitCode              *int
	DurationMS            *int
	TeacherComment        string
	ReviewedBy            *uuid.UUID
	ReviewedAt            *time.Time
	AttemptNumber         int
	PracticeTitle         string
	StudentEmail          string
	CourseTitle           string
	LessonTitle           string
	ProgressStatus        string
	ProgressStartedAt     *time.Time
	ProgressCompletedAt   *time.Time
	ProgressLastAttemptAt *time.Time
	ProgressAttemptsCount int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type CreateSubmissionRequest struct {
	PracticeID uuid.UUID
	StudentID  uuid.UUID
	Code       string
	Language   string
	Output     string
	Error      string
	ErrorType  string
	ExitCode   *int
	DurationMS *int
}

type ReviewSubmissionRequest struct {
	SubmissionID uuid.UUID
	TeacherID    uuid.UUID
	IsAdmin      bool
	Status       string
	Comment      string
}

type StudentListFilter struct {
	StudentID  uuid.UUID
	CourseID   *uuid.UUID
	PracticeID *uuid.UUID
	Status     string
}

type TeacherListFilter struct {
	TeacherID  uuid.UUID
	IsAdmin    bool
	CourseID   *uuid.UUID
	PracticeID *uuid.UUID
	StudentID  *uuid.UUID
	Status     string
	Limit      int
	Offset     int
}
