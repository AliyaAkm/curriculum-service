package reviewcourse

import (
	"github.com/google/uuid"
	"time"
)

type ReviewCourseRequest struct {
	IsApproved bool      `json:"is_approved"`
	AdminID    uuid.UUID `json:"admin_id"`
	Comment    string    `json:"comment"`
}

type ReviewCourseResponse struct {
	CourseID   uuid.UUID `json:"course_id"`
	IsChecked  bool      `json:"is_checked"`
	IsApproved bool      `json:"is_approved"`
	Comment    string    `json:"comment"`
	ReviewedBy uuid.UUID `json:"reviewed_by"`
	ReviewedAt time.Time `json:"reviewed_at"`
}
