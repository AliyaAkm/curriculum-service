package title

import (
	"curriculum-service/internal/domain/locale"
	"github.com/google/uuid"
)

type LessonTitleModel struct {
	ID       uuid.UUID `gorm:"column:id;primary_key"`
	Name     string    `gorm:"column:name"`
	LocaleID uuid.UUID `gorm:"column:locale_id;foreign_key"`
	LessonID uuid.UUID `gorm:"column:lesson_id;foreign_key"`
	Locale   locale.Locale
}

func (LessonTitleModel) TableName() string {
	return "course_lesson_titles"
}
