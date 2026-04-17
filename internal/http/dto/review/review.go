package review

import (
	"github.com/google/uuid"
	"time"
)

type UpdateReviewRequest struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}

type ReviewRequest struct {
	CourseID uuid.UUID `json:"course_id" validate:"required"`
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	Rating   int       `json:"rating" validate:"required,gte=1,lte=5"`
	Comment  string    `json:"comment"`
}

type ReviewResponse struct {
	ID        uuid.UUID `json:"id"`
	CourseID  uuid.UUID `json:"course_id"`
	UserID    uuid.UUID `json:"user_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	ViewCount int       `json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
