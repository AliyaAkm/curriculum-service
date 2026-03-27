package lesson

import (
	"curriculum-service/internal/domain/keypoint"
	"curriculum-service/internal/domain/outcome"
	"curriculum-service/internal/domain/summary"
	"curriculum-service/internal/domain/theorycontent"
	"curriculum-service/internal/domain/title"
	"github.com/google/uuid"
	"time"
)

type LessonModel struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	ModuleID        uuid.UUID `gorm:"column:module_id;type:uuid"`
	DurationMinutes int       `gorm:"column:duration_minutes"`
	XPReward        int       `gorm:"column:xp_reward"`
	CodeSnippet     *string   `gorm:"column:code_snippet"`
	ExampleOutput   *string   `gorm:"column:example_output"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`

	Titles         []title.LessonTitleModel                 `gorm:"foreignKey:LessonID;references:ID"`
	Summaries      []summary.LessonSummaryModel             `gorm:"foreignKey:LessonID;references:ID"`
	Outcomes       []outcome.LessonOutcomeModel             `gorm:"foreignKey:LessonID;references:ID"`
	TheoryContents []theorycontent.LessonTheoryContentModel `gorm:"foreignKey:LessonID;references:ID"`
	KeyPoints      []keypoint.LessonKeyPointModel           `gorm:"foreignKey:LessonID;references:ID"`
}

func (LessonModel) TableName() string {
	return "course_lessons"
}
