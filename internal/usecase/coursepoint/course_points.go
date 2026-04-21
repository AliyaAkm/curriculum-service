package coursepoint

import (
	"context"
	"curriculum-service/internal/domain/coursepoint"
	"github.com/google/uuid"
)

func (u *UseCase) CreateCoursePoint(ctx context.Context, value *coursepoint.UserCoursePoints) (*coursepoint.UserCoursePoints, error) {
	value.ID = uuid.New()

	err := u.repo.CreateCoursePoint(ctx, value)
	if err != nil {
		return nil, err
	}
	return u.repo.GetCoursePointByID(ctx, value.ID)
}

func (u *UseCase) UpdateCoursePoint(ctx context.Context, id uuid.UUID, newValue *coursepoint.UserCoursePoints) (*coursepoint.UserCoursePoints, error) {
	newValue.ID = id

	oldValue, err := u.repo.GetCoursePointByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if newValue.XP != 0 {
		oldValue.XP = newValue.XP
	}
	err = u.repo.UpdateCoursePoint(ctx, id, newValue)
	if err != nil {
		return nil, err
	}
	return u.repo.GetCoursePointByID(ctx, id)
}

func (u *UseCase) DeleteCoursePoint(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteCoursePoint(ctx, id)
}
func (u *UseCase) GetCoursePointByCourseID(ctx context.Context, id uuid.UUID) ([]coursepoint.Leaderboard, error) {
	value, err := u.repo.GetCoursePointByCourseID(ctx, id)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(value); i++ {
		value[i].Place = i + 1
	}
	return value, nil
}
