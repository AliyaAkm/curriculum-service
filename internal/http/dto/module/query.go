package module

import (
	"github.com/google/uuid"
	"time"
)

type ModuleRequest struct {
	CourseID uuid.UUID `json:"course_id"`
	Title    string    `json:"title"`
	Summary  string    `json:"summary"`
	Locale   string    `json:"locale"`
}

type GetModuleQuery struct {
	CourseID uuid.UUID `form:"course_id"`
	Locale   string    `form:"locale"`
	Page     int       `form:"page"`
	Limit    int       `form:"limit"`
}

type Modules struct {
	ID        uuid.UUID `json:"id"`
	CourseID  uuid.UUID `json:"course_id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Locale    string    `json:"locale"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
