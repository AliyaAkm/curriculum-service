package practice

import (
	"context"
	practicedomain "curriculum-service/internal/domain/practice"

	"github.com/google/uuid"
)

type client interface {
	Create(ctx context.Context, value practicedomain.Task) (*practicedomain.Task, error)
	Update(ctx context.Context, id uuid.UUID, value practicedomain.TaskUpdate) (*practicedomain.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*practicedomain.Task, error)
	ListByLesson(ctx context.Context, lessonID uuid.UUID) ([]practicedomain.Task, error)
}
