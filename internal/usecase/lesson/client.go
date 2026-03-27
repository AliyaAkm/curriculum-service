package lesson

import (
	"context"
	lessondomain "curriculum-service/internal/domain/lesson"
	"github.com/google/uuid"
)

type Repository interface {
	GetAllLessons(ctx context.Context, moduleID uuid.UUID) ([]lessondomain.LessonModel, error)
	GetLessonByID(ctx context.Context, id uuid.UUID) (*lessondomain.LessonModel, error)
	CreateLesson(ctx context.Context, value *lessondomain.LessonModel) error
	UpdateLesson(ctx context.Context, id uuid.UUID, value *lessondomain.LessonModel) error
}
