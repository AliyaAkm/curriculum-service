package course

import (
	"context"
	"curriculum-service/internal/domain/course"
	dtocourse "curriculum-service/internal/http/dto/course"
	"github.com/google/uuid"
)

type Repository interface {
	GetAllCourses(ctx context.Context, query dtocourse.GetCoursesQuery) ([]course.Course, error)
	CreateCourse(ctx context.Context, value *course.Course) (uuid.UUID, error)
	GetCourseByID(ctx context.Context, id uuid.UUID) (*course.Course, error)
	DeleteCourse(ctx context.Context, id uuid.UUID) error
	UpdateCourse(ctx context.Context, id uuid.UUID, value *course.Course) error
}
