package coursepoint

import (
	"context"
	"curriculum-service/internal/domain/coursepoint"
	"github.com/google/uuid"
)

type client interface {
	CreateCoursePoint(ctx context.Context, value *coursepoint.UserCoursePoints) (*coursepoint.UserCoursePoints, error)
	UpdateCoursePoint(ctx context.Context, id uuid.UUID, newValue *coursepoint.UserCoursePoints) (*coursepoint.UserCoursePoints, error)
	DeleteCoursePoint(ctx context.Context, id uuid.UUID) error
	GetCoursePointByCourseID(ctx context.Context, id uuid.UUID) ([]coursepoint.Leaderboard, error)
}
