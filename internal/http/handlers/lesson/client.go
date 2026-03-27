package lesson

import (
	"context"
	lessondomain "curriculum-service/internal/domain/lesson"
	"curriculum-service/internal/domain/locale"
	"github.com/google/uuid"
)

type client interface {
	GetAllLessons(ctx context.Context, moduleID uuid.UUID) ([]lessondomain.LessonModel, error)
	CreateLesson(ctx context.Context, value *lessondomain.LessonModel) (*lessondomain.LessonModel, error)
	UpdateLesson(ctx context.Context, id uuid.UUID, value *lessondomain.LessonModel) (*lessondomain.LessonModel, error)
	GetLessonByID(ctx context.Context, id uuid.UUID) (*lessondomain.LessonModel, error)
}

type localClient interface {
	GetAllLocales(ctx context.Context) ([]locale.Locale, error)
}
