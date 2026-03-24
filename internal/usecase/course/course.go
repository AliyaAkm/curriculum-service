package course

import (
	"context"
	"curriculum-service/internal/domain/course"
	dtocourse "curriculum-service/internal/http/dto/course"
)

func (u *UseCase) GetAllCourses(ctx context.Context, query dtocourse.GetCoursesQuery) ([]course.Course, error) {
	resp, err := u.repo.GetAllCourses(ctx, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (u *UseCase) CreateCourse(ctx context.Context, value *course.Course) (*course.Course, error) {
	return u.repo.CreateCourse(ctx, value)
}
