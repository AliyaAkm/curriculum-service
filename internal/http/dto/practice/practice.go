package practice

import (
	"time"

	"github.com/google/uuid"
)

type Practice struct {
	ID               uuid.UUID `json:"id"`
	LessonID         uuid.UUID `json:"lesson_id"`
	Position         int       `json:"position"`
	Title            Locale    `json:"title"`
	Summary          Locale    `json:"summary"`
	Brief            Locale    `json:"brief"`
	StarterCode      string    `json:"starter_code"`
	SuccessCriteria  []Locale  `json:"success_criteria"`
	KnowledgeChecks  []Locale  `json:"knowledge_checks"`
	PromptSuggestion Locale    `json:"prompt_suggestion"`
	XPReward         int       `json:"xp_reward"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PracticeRequest struct {
	LessonID         uuid.UUID `json:"lesson_id"`
	Position         int       `json:"position"`
	Title            Locale    `json:"title"`
	Summary          Locale    `json:"summary"`
	Brief            Locale    `json:"brief"`
	StarterCode      string    `json:"starter_code"`
	SuccessCriteria  []Locale  `json:"success_criteria"`
	KnowledgeChecks  []Locale  `json:"knowledge_checks"`
	PromptSuggestion Locale    `json:"prompt_suggestion"`
	XPReward         int       `json:"xp_reward"`
}

type Locale struct {
	EN string `json:"en"`
	RU string `json:"ru"`
	KK string `json:"kk"`
}
