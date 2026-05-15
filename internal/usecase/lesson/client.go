package lesson

import (
	"context"
	lessondomain "curriculum-service/internal/domain/lesson"
	"github.com/google/uuid"
)

type Repository interface {
	GetAllLessons(ctx context.Context, moduleID uuid.UUID) ([]lessondomain.LessonModel, error)
	GetLessonByID(ctx context.Context, id uuid.UUID) (*lessondomain.LessonModel, error)
	GetCourseIDByModuleID(ctx context.Context, moduleID uuid.UUID) (uuid.UUID, error)
	GetLessonAccessInfo(ctx context.Context, lessonID uuid.UUID) (uuid.UUID, uuid.UUID, int, error)
	HasSubscription(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (bool, error)
	IsModuleInFreePreview(ctx context.Context, courseID uuid.UUID, moduleID uuid.UUID, limit int) (bool, error)
	CreateLesson(ctx context.Context, value *lessondomain.LessonModel) error
	UpdateLesson(ctx context.Context, id uuid.UUID, value *lessondomain.LessonModel) error
	UpdateLessonVideoObjectKey(ctx context.Context, id uuid.UUID, videoObjectKey *string) error
}
