package course

import (
	"context"
	"curriculum-service/internal/domain/course"
)

func (u *UseCase) GetAllCourses(ctx context.Context) ([]course.Courses, error) {
	return u.repo.GetAllCourses(ctx)
}
