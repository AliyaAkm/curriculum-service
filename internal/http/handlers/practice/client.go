package practice

import (
	"context"
	practicedomain "curriculum-service/internal/domain/practice"

	"github.com/google/uuid"
)

type client interface {
	CreatePractice(ctx context.Context, value *practicedomain.Practice) (*practicedomain.Practice, error)
	GetPracticeByLessonIDForUser(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID, hasFullAccess bool) (*practicedomain.Practice, error)
	GetPracticeByIDForUser(ctx context.Context, userID uuid.UUID, id uuid.UUID, hasFullAccess bool) (*practicedomain.Practice, error)
	UpdatePractice(ctx context.Context, id uuid.UUID, value *practicedomain.Practice) (*practicedomain.Practice, error)
	DeletePractice(ctx context.Context, id uuid.UUID) error
}
