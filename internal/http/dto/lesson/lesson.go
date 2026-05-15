package lesson

import (
	"github.com/google/uuid"
	"time"
)

type Lesson struct { // course_lessons
	ID       uuid.UUID `json:"id"`
	ModuleID uuid.UUID `json:"module_id"`
	Position int       `json:"position"`

	Title         Locale   `json:"title"`
	Summary       Locale   `json:"summary"`
	OutCome       Locale   `json:"out_come"`
	KeyPoints     []Locale `json:"key_points"`
	TheoryContent Locale   `json:"theory_content"`

	DurationMinutes int       `json:"duration_minutes"`
	XPReward        int       `json:"xp_reward"`
	CodeSnippet     *string   `json:"code_snippet"`
	ExampleOutput   *string   `json:"example_output"`
	VideoObjectKey  *string   `json:"video_object_key,omitempty"`
	VideoURL        *string   `json:"video_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Locale struct {
	EN string `json:"en"`
	KK string `json:"kk"`
	RU string `json:"ru"`
}
