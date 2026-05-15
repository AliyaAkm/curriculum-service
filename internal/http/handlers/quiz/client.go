package quiz

import (
	"context"
	quizdomain "curriculum-service/internal/domain/quiz"

	"github.com/google/uuid"
)

type client interface {
	CreateQuiz(ctx context.Context, value *quizdomain.Quiz) (*quizdomain.Quiz, error)
	GetQuizzesByLessonIDForUser(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID, hasFullAccess bool) ([]quizdomain.Quiz, error)
	GetQuizByIDForUser(ctx context.Context, userID uuid.UUID, id uuid.UUID, hasFullAccess bool) (*quizdomain.Quiz, error)
	SubmitAnswer(ctx context.Context, userID uuid.UUID, quizID uuid.UUID, selectedAnswerIndex int, hasFullAccess bool) (*quizdomain.AnswerResult, error)
	UpdateQuiz(ctx context.Context, id uuid.UUID, value *quizdomain.Quiz) (*quizdomain.Quiz, error)
	DeleteQuiz(ctx context.Context, id uuid.UUID) error
}
