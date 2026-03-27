package lesson

import (
	"curriculum-service/internal/domain/keypoint"
	"curriculum-service/internal/domain/outcome"
	"curriculum-service/internal/domain/summary"
	"curriculum-service/internal/domain/theorycontent"
	"curriculum-service/internal/domain/title"
	"github.com/google/uuid"
)

type LessonRequest struct {
	ModuleID        uuid.UUID `json:"module_id"`
	DurationMinutes int       `json:"duration_minutes"`
	XPReward        int       `json:"xp_reward"`
	CodeSnippet     string    `json:"code_snippet"`
	ExampleOutput   string    `json:"example_output"`
	Title           Locale    `json:"title"`
	Summary         Locale    `json:"summary"`
	OutCome         Locale    `json:"outcome"`
	TheoryContent   Locale    `json:"theory_content"`
	KeyPoints       []Locale  `json:"key_points"`
}

type CreateLessonRequest struct {
	Titles          []title.LessonTitleModel                 `json:"titles"`
	Summaries       []summary.LessonSummaryModel             `json:"summaries"`
	Outcomes        []outcome.LessonOutcomeModel             `json:"outcomes"`
	TheoryContents  []theorycontent.LessonTheoryContentModel `json:"theory_contents"`
	KeyPoints       []keypoint.LessonKeyPointModel           `json:"key_points"`
	DurationMinutes int                                      `json:"duration_minutes"`
	XPReward        int                                      `json:"xp_reward"`
	CodeSnippet     string                                   `json:"code_snippet"`
	ExampleOutput   string                                   `json:"example_output"`
}

type UpdateLessonRequest struct {
	Titles          []title.LessonTitleModel                 `json:"titles"`
	Summaries       []summary.LessonSummaryModel             `json:"summaries"`
	Outcomes        []outcome.LessonOutcomeModel             `json:"outcomes"`
	TheoryContents  []theorycontent.LessonTheoryContentModel `json:"theory_contents"`
	KeyPoints       []keypoint.LessonKeyPointModel           `json:"key_points"`
	DurationMinutes int                                      `json:"duration_minutes"`
	XPReward        int                                      `json:"xp_reward"`
	CodeSnippet     string                                   `json:"code_snippet"`
	ExampleOutput   string                                   `json:"example_output"`
}
