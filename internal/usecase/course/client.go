package course

import (
	"context"
	"curriculum-service/internal/domain/course"
)

type Repository interface {
	GetAllCourses(ctx context.Context) ([]course.Course, error)
}
