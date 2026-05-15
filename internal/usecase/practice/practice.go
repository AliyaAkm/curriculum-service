package practice

import (
	"context"
	"curriculum-service/internal/domain"
	practicedomain "curriculum-service/internal/domain/practice"

	"github.com/google/uuid"
)

const freePreviewModulesLimit = 2

func (u *UseCase) CreatePractice(ctx context.Context, value *practicedomain.Practice) (*practicedomain.Practice, error) {
	value.ID = uuid.New()

	if err := u.repo.CreatePractice(ctx, value); err != nil {
		return nil, err
	}

	return u.repo.GetPracticeByID(ctx, value.ID)
}

func (u *UseCase) GetPracticeByLessonIDForUser(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID, hasFullAccess bool) (*practicedomain.Practice, error) {
	if err := u.requireLessonAccess(ctx, userID, lessonID, hasFullAccess); err != nil {
		return nil, err
	}

	return u.repo.GetPracticeByLessonID(ctx, lessonID)
}

func (u *UseCase) GetPracticeByIDForUser(ctx context.Context, userID uuid.UUID, id uuid.UUID, hasFullAccess bool) (*practicedomain.Practice, error) {
	practiceEntity, err := u.repo.GetPracticeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := u.requireLessonAccess(ctx, userID, practiceEntity.LessonID, hasFullAccess); err != nil {
		return nil, err
	}

	return practiceEntity, nil
}

func (u *UseCase) UpdatePractice(ctx context.Context, id uuid.UUID, value *practicedomain.Practice) (*practicedomain.Practice, error) {
	if err := u.repo.UpdatePractice(ctx, id, value); err != nil {
		return nil, err
	}

	return u.repo.GetPracticeByID(ctx, id)
}

func (u *UseCase) DeletePractice(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeletePractice(ctx, id)
}

func (u *UseCase) requireLessonAccess(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID, hasFullAccess bool) error {
	if hasFullAccess {
		return nil
	}

	courseID, moduleID, _, err := u.repo.GetLessonAccessInfo(ctx, lessonID)
	if err != nil {
		return err
	}

	hasSubscription, err := u.repo.HasSubscription(ctx, userID, courseID)
	if err != nil {
		return err
	}
	if hasSubscription {
		return nil
	}

	allowed, err := u.repo.IsModuleInFreePreview(ctx, courseID, moduleID, freePreviewModulesLimit)
	if err != nil {
		return err
	}
	if !allowed {
		return domain.ErrForbidden
	}

	return nil
}
