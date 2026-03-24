package course

import (
	"context"
	"curriculum-service/internal/domain/course"
)

type client interface {
	GetAllCourses(ctx context.Context) ([]course.Course, error)
}
