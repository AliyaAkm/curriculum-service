package course

import (
	"context"
	"curriculum-service/internal/domain/course"
)

func (r *Repo) GetAllCourses(ctx context.Context) ([]course.Courses, error) {
	var courses []course.Courses
	err := r.db.WithContext(ctx).
		Preload("Status").
		Preload("Level").
		Preload("DurationCategory").
		Preload("Author.Roles").
		Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
