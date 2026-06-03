package lesson

import (
	"context"
	lessondomain "curriculum-service/internal/domain/lesson"
	"curriculum-service/internal/domain/locale"
	"github.com/google/uuid"
	"io"
)

type client interface {
	GetAllLessons(ctx context.Context, moduleID uuid.UUID) ([]lessondomain.LessonModel, error)
	GetAllLessonsForUser(ctx context.Context, userID uuid.UUID, moduleID uuid.UUID, hasFullAccess bool) ([]lessondomain.LessonModel, error)
	CreateLesson(ctx context.Context, value *lessondomain.LessonModel) (*lessondomain.LessonModel, error)
	UpdateLesson(ctx context.Context, id uuid.UUID, value *lessondomain.LessonModel) (*lessondomain.LessonModel, error)
	UpdateLessonVideoObjectKey(ctx context.Context, id uuid.UUID, videoObjectKey *string) (*lessondomain.LessonModel, error)
	GetLessonByID(ctx context.Context, id uuid.UUID) (*lessondomain.LessonModel, error)
	GetLessonByIDForUser(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID, hasFullAccess bool) (*lessondomain.LessonModel, error)
}

type localClient interface {
	GetAllLocales(ctx context.Context) ([]locale.Locale, error)
}

type objectStorage interface {
	PutObject(ctx context.Context, objectKey string, body io.ReadSeeker, size int64, contentType string) error
	PresignGetObject(objectKey string) (string, error)
}
