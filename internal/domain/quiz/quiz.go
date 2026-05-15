package quiz

import (
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	ID                 uuid.UUID
	LessonID           uuid.UUID
	Position           int
	Question           Locale
	Options            []Option
	CorrectAnswerIndex int
	Explanation        Locale
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Option struct {
	ID       uuid.UUID
	Position int
	Text     Locale
}

type AnswerResult struct {
	QuizID              uuid.UUID
	SelectedAnswerIndex int
	IsCorrect           bool
	CorrectAnswerIndex  int
	CorrectOptionID     uuid.UUID
	Explanation         Locale
}

type Locale struct {
	EN string
	RU string
	KK string
}
