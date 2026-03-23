package tag

import "github.com/google/uuid"

type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}
