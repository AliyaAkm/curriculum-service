package practice

import (
	"context"
	"curriculum-service/internal/domain/practice"

	"github.com/google/uuid"
)

type Repository interface {
	CreatePractice(ctx context.Context, value *practice.Practice) error
	GetPracticeByID(ctx context.Context, id uuid.UUID) (*practice.Practice, error)
	GetPracticeByLessonID(ctx context.Context, lessonID uuid.UUID) (*practice.Practice, error)
	UpdatePractice(ctx context.Context, id uuid.UUID, value *practice.Practice) error
	DeletePractice(ctx context.Context, id uuid.UUID) error
	GetLessonAccessInfo(ctx context.Context, lessonID uuid.UUID) (uuid.UUID, uuid.UUID, int, error)
	HasSubscription(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (bool, error)
	IsModuleInFreePreview(ctx context.Context, courseID uuid.UUID, moduleID uuid.UUID, limit int) (bool, error)
}
