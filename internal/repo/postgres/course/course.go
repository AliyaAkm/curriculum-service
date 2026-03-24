package course

import (
	"context"
	"curriculum-service/internal/domain/course"
)

func (r *Repo) GetAllCourses(ctx context.Context) ([]course.Course, error) {
	var courses []course.Course
	err := r.db.WithContext(ctx).
		Preload("Status").
		Preload("Level").
		Preload("DurationCategory").
		Preload("Author.Roles").
		Preload("Tags").
		Preload("Topic").
		Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
