package status

import (
	"context"
	"curriculum-service/internal/domain/status"
)

func (r *Repo) GetAllStatus(ctx context.Context) ([]status.Status, error) {
	var statuses []status.Status
	// select * from course_statuses -> column:id ->id (с бд значение)
	err := r.db.WithContext(ctx).Order("code ASC").Find(&statuses).Error
	if err != nil {
		return nil, err
	}

	return statuses, nil
}
