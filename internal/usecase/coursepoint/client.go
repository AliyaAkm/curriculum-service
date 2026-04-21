package coursepoint

import (
	"context"
	"curriculum-service/internal/domain/coursepoint"
	"github.com/google/uuid"
)

type Repository interface {
	CreateCoursePoint(ctx context.Context, value *coursepoint.UserCoursePoints) error
	GetCoursePointByID(ctx context.Context, id uuid.UUID) (*coursepoint.UserCoursePoints, error)
	UpdateCoursePoint(ctx context.Context, id uuid.UUID, value *coursepoint.UserCoursePoints) error
	DeleteCoursePoint(ctx context.Context, id uuid.UUID) error
	GetCoursePointByCourseID(ctx context.Context, id uuid.UUID) ([]coursepoint.Leaderboard, error )
}
