package course

import (
	"context"
	"curriculum-service/internal/domain/course"
	dtocourse "curriculum-service/internal/http/dto/course"
)

type Repository interface {
	GetAllCourses(ctx context.Context, query dtocourse.GetCoursesQuery) ([]course.Course, error)
	CreateCourse(ctx context.Context, value *course.Course) (*course.Course, error)
}
