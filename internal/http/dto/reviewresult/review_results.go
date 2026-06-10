package reviewresult

import (
	"github.com/google/uuid"
	"time"
)

type CourseReviewResponse struct {
	CourseID     uuid.UUID  `json:"course_id"`
	Title        string     `json:"title"`
	IsChecked    bool       `json:"is_checked"`
	IsApproved   bool       `json:"is_approved"`
	ReviewStatus string     `json:"review_status"`
	Comment      *string    `json:"comment,omitempty"`
	ReviewedBy   *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`
}
