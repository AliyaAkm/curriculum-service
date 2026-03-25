package module

import "github.com/google/uuid"

type Module struct {
	ID          uuid.UUID `json:"id"`
	CourseID    uuid.UUID `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Locale      string    `json:"locale"`
}
