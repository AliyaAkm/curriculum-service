package resubmit

import "github.com/google/uuid"

type ResubmitReviewResponse struct {
	CourseID   uuid.UUID `json:"course_id"`
	IsChecked  bool      `json:"is_checked"`
	IsApproved bool      `json:"is_approved"`
	Message    string    `json:"message"`
}
