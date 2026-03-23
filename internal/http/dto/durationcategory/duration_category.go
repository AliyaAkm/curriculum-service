package durationcategory

import "github.com/google/uuid"

type DurationCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
}
