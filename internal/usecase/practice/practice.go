package practice

import (
	"context"
	"curriculum-service/internal/domain"
	practicedomain "curriculum-service/internal/domain/practice"
	"strings"

	"github.com/google/uuid"
)

func (u *UseCase) Create(ctx context.Context, value practicedomain.Task) (*practicedomain.Task, error) {
	if value.XPReward == 0 {
		value.XPReward = 25
	}
	if strings.TrimSpace(value.CheckType) == "" {
		value.CheckType = practicedomain.CheckTypeAuto
	}
	value.CheckType = strings.TrimSpace(value.CheckType)
	if err := validate(value); err != nil {
		return nil, err
	}
	value.ID = uuid.New()
	if err := u.repo.Create(ctx, value); err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, value.ID)
}

func (u *UseCase) Update(ctx context.Context, id uuid.UUID, value practicedomain.TaskUpdate) (*practicedomain.Task, error) {
	if id == uuid.Nil {
		return nil, domain.ErrValidation
	}
	if err := validateUpdate(value); err != nil {
		return nil, err
	}
	if value.CheckType != nil {
		trimmed := strings.TrimSpace(*value.CheckType)
		value.CheckType = &trimmed
	}
	if err := u.repo.Update(ctx, id, value); err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return domain.ErrValidation
	}
	return u.repo.Delete(ctx, id)
}

func (u *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*practicedomain.Task, error) {
	if id == uuid.Nil {
		return nil, domain.ErrValidation
	}
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) ListByLesson(ctx context.Context, lessonID uuid.UUID) ([]practicedomain.Task, error) {
	if lessonID == uuid.Nil {
		return nil, domain.ErrValidation
	}
	return u.repo.ListByLesson(ctx, lessonID)
}

func validate(value practicedomain.Task) error {
	if value.LessonID == uuid.Nil ||
		value.Position <= 0 ||
		strings.TrimSpace(value.Title) == "" ||
		strings.TrimSpace(value.Language) == "" ||
		strings.TrimSpace(value.StarterCode) == "" ||
		value.XPReward <= 0 ||
		!isCheckType(value.CheckType) {
		return domain.ErrValidation
	}
	if value.CheckType == practicedomain.CheckTypeAuto && strings.TrimSpace(value.ExpectedOutput) == "" {
		return domain.ErrValidation
	}
	return nil
}

func validateUpdate(value practicedomain.TaskUpdate) error {
	if value.Title == nil &&
		value.Description == nil &&
		value.Language == nil &&
		value.StarterCode == nil &&
		value.ExpectedOutput == nil &&
		value.XPReward == nil &&
		value.CheckType == nil {
		return domain.ErrValidation
	}
	if value.XPReward != nil && *value.XPReward <= 0 {
		return domain.ErrValidation
	}
	if value.Title != nil && strings.TrimSpace(*value.Title) == "" {
		return domain.ErrValidation
	}
	if value.Language != nil && strings.TrimSpace(*value.Language) == "" {
		return domain.ErrValidation
	}
	if value.StarterCode != nil && strings.TrimSpace(*value.StarterCode) == "" {
		return domain.ErrValidation
	}
	if value.CheckType != nil && !isCheckType(*value.CheckType) {
		return domain.ErrValidation
	}
	return nil
}

func isCheckType(value string) bool {
	switch strings.TrimSpace(value) {
	case practicedomain.CheckTypeAuto, practicedomain.CheckTypeManual:
		return true
	default:
		return false
	}
}
