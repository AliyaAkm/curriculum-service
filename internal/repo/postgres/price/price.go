package price

import (
	"context"
	"curriculum-service/internal/domain/price"
	"github.com/google/uuid"
)

func (r *Repo) GetPriceByCourseID(ctx context.Context, courseID uuid.UUID) (*price.CoursePrice, error) {
	var entity *price.CoursePrice
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
