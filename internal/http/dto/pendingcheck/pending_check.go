package pendingcheck

import (
	"github.com/google/uuid"
	"time"
)

type PendingCheckCourseResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AuthorID    uuid.UUID `json:"author_id"`
	Category    string    `json:"category"`   // empty if topic not set
	Difficulty  string    `json:"difficulty"` // empty if level not set
	IsChecked   bool      `json:"is_checked"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
