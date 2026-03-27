package lesson

import (
	"context"
	lessondomain "curriculum-service/internal/domain/lesson"
	"github.com/google/uuid"
)

func (u *UseCase) GetAllLessons(ctx context.Context, moduleID uuid.UUID) ([]lessondomain.LessonModel, error) {
	return u.repo.GetAllLessons(ctx, moduleID)
}

func (u *UseCase) GetLessonByID(ctx context.Context, id uuid.UUID) (*lessondomain.LessonModel, error) {
	return u.repo.GetLessonByID(ctx, id)
}

func (u *UseCase) CreateLesson(ctx context.Context, value *lessondomain.LessonModel) (*lessondomain.LessonModel, error) {
	value.ID = uuid.New()

	err := u.repo.CreateLesson(ctx, value)
	if err != nil {
		return nil, err
	}

	return u.repo.GetLessonByID(ctx, value.ID)
}

func (u *UseCase) UpdateLesson(ctx context.Context, id uuid.UUID, value *lessondomain.LessonModel) (*lessondomain.LessonModel, error) {
	value.ID = id

	if err := u.repo.UpdateLesson(ctx, id, value); err != nil {
		return nil, err
	}

	return u.repo.GetLessonByID(ctx, id)
}
