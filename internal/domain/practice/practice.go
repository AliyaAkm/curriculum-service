package practice

import (
	"time"

	"github.com/google/uuid"
)

type Practice struct {
	ID               uuid.UUID
	LessonID         uuid.UUID
	Position         int
	Title            Locale
	Summary          Locale
	Brief            Locale
	StarterCode      string
	SuccessCriteria  []Locale
	KnowledgeChecks  []Locale
	PromptSuggestion Locale
	XPReward         int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Locale struct {
	EN string `json:"en"`
	RU string `json:"ru"`
	KK string `json:"kk"`
}
