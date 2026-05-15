package quiz

import (
	"context"
	quizdomain "curriculum-service/internal/domain/quiz"

	"github.com/google/uuid"
)

type Repository interface {
	CreateQuiz(ctx context.Context, value *quizdomain.Quiz) error
	GetQuizByID(ctx context.Context, id uuid.UUID) (*quizdomain.Quiz, error)
	GetQuizzesByLessonID(ctx context.Context, lessonID uuid.UUID) ([]quizdomain.Quiz, error)
	UpdateQuiz(ctx context.Context, id uuid.UUID, value *quizdomain.Quiz) error
	DeleteQuiz(ctx context.Context, id uuid.UUID) error
	SaveQuizAttempt(ctx context.Context, userID uuid.UUID, quizID uuid.UUID, selectedAnswerIndex int, isCorrect bool) error
	GetLessonAccessInfo(ctx context.Context, lessonID uuid.UUID) (uuid.UUID, uuid.UUID, int, error)
	HasSubscription(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (bool, error)
	IsModuleInFreePreview(ctx context.Context, courseID uuid.UUID, moduleID uuid.UUID, limit int) (bool, error)
}
