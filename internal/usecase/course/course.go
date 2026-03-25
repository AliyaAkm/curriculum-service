package course

import (
	"context"
	"curriculum-service/internal/domain/course"
	dtocourse "curriculum-service/internal/http/dto/course"
	"github.com/google/uuid"
)

func (u *UseCase) GetAllCourses(ctx context.Context, query dtocourse.GetCoursesQuery) ([]course.Course, error) {
	resp, err := u.repo.GetAllCourses(ctx, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (u *UseCase) CreateCourse(ctx context.Context, value *course.Course) (*course.Course, error) {
	value.ID = uuid.New()
	id, err := u.repo.CreateCourse(ctx, value)
	if err != nil {
		return nil, err
	}
	return u.repo.GetCourseByID(ctx, id)
}

func (u *UseCase) GetCourseByID(ctx context.Context, id uuid.UUID) (*course.Course, error) {
	return u.repo.GetCourseByID(ctx, id)
}

func (u *UseCase) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteCourse(ctx, id)
}
func (u *UseCase) UpdateCourse(ctx context.Context, id uuid.UUID, value *course.Course) (*course.Course, error) {
	err := u.repo.UpdateCourse(ctx, id, value)
	if err != nil {
		return nil, err
	}
	return u.repo.GetCourseByID(ctx, id)
}
