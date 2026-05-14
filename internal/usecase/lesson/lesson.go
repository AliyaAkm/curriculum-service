package lesson

import (
	"context"
	"curriculum-service/internal/domain"
	lessondomain "curriculum-service/internal/domain/lesson"
	"github.com/google/uuid"
)

const freePreviewModulesLimit = 2

func (u *UseCase) GetAllLessons(ctx context.Context, moduleID uuid.UUID) ([]lessondomain.LessonModel, error) {
	return u.repo.GetAllLessons(ctx, moduleID)
}

func (u *UseCase) GetLessonByID(ctx context.Context, id uuid.UUID) (*lessondomain.LessonModel, error) {
	return u.repo.GetLessonByID(ctx, id)
}

func (u *UseCase) GetAllLessonsForUser(ctx context.Context, userID uuid.UUID, moduleID uuid.UUID, hasFullAccess bool) ([]lessondomain.LessonModel, error) {
	if !hasFullAccess {
		courseID, err := u.repo.GetCourseIDByModuleID(ctx, moduleID)
		if err != nil {
			return nil, err
		}

		hasSubscription, err := u.repo.HasSubscription(ctx, userID, courseID)
		if err != nil {
			return nil, err
		}

		if !hasSubscription {
			allowed, err := u.repo.IsModuleInFreePreview(ctx, courseID, moduleID, freePreviewModulesLimit)
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, domain.ErrForbidden
			}
		}
	}

	return u.repo.GetAllLessons(ctx, moduleID)
}

func (u *UseCase) GetLessonByIDForUser(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID, hasFullAccess bool) (*lessondomain.LessonModel, error) {
	if !hasFullAccess {
		courseID, moduleID, _, err := u.repo.GetLessonAccessInfo(ctx, lessonID)
		if err != nil {
			return nil, err
		}

		hasSubscription, err := u.repo.HasSubscription(ctx, userID, courseID)
		if err != nil {
			return nil, err
		}

		if !hasSubscription {
			allowed, err := u.repo.IsModuleInFreePreview(ctx, courseID, moduleID, freePreviewModulesLimit)
			if err != nil {
				return nil, err
			}
			if !allowed {
				return nil, domain.ErrForbidden
			}
		}
	}

	return u.repo.GetLessonByID(ctx, lessonID)
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
