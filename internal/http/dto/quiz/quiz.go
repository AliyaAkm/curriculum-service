package quiz

import (
	"time"

	"github.com/google/uuid"
)

type Quiz struct {
	ID                 uuid.UUID `json:"id"`
	LessonID           uuid.UUID `json:"lesson_id"`
	Position           int       `json:"position"`
	Question           Locale    `json:"question"`
	Options            []Option  `json:"options"`
	CorrectAnswerIndex int       `json:"correct_answer_index"`
	CorrectOptionID    uuid.UUID `json:"correct_option_id"`
	Explanation        Locale    `json:"explanation"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Option struct {
	ID       uuid.UUID `json:"id"`
	Position int       `json:"position"`
	Text     Locale    `json:"text"`
}

type QuizRequest struct {
	LessonID           uuid.UUID `json:"lesson_id"`
	Position           int       `json:"position"`
	Question           Locale    `json:"question"`
	Options            []Locale  `json:"options"`
	CorrectAnswerIndex int       `json:"correct_answer_index"`
	Explanation        Locale    `json:"explanation"`
}

type QuizUpdateRequest struct {
	Position           int      `json:"position"`
	Question           Locale   `json:"question"`
	Options            []Locale `json:"options"`
	CorrectAnswerIndex int      `json:"correct_answer_index"`
	Explanation        Locale   `json:"explanation"`
}

type AnswerRequest struct {
	SelectedAnswerIndex int `json:"selected_answer_index"`
}

type AnswerResponse struct {
	QuizID              uuid.UUID `json:"quiz_id"`
	SelectedAnswerIndex int       `json:"selected_answer_index"`
	IsCorrect           bool      `json:"is_correct"`
	CorrectAnswerIndex  int       `json:"correct_answer_index"`
	CorrectOptionID     uuid.UUID `json:"correct_option_id"`
	Explanation         Locale    `json:"explanation"`
}

type Locale struct {
	EN string `json:"en"`
	RU string `json:"ru"`
	KK string `json:"kk"`
}
