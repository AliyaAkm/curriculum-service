package coursepoint

import (
	"context"
	"curriculum-service/internal/domain/coursepoint"
	"github.com/google/uuid"
)

func (r *Repo) CreateCoursePoint(ctx context.Context, value *coursepoint.UserCoursePoints) error {
	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetCoursePointByID(ctx context.Context, id uuid.UUID) (*coursepoint.UserCoursePoints, error) {
	var entity coursepoint.UserCoursePoints
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repo) UpdateCoursePoint(ctx context.Context, id uuid.UUID, value *coursepoint.UserCoursePoints) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(&value).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteCoursePoint(ctx context.Context, id uuid.UUID) error {
	var entity coursepoint.UserCoursePoints
	err := r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetCoursePointByCourseID(ctx context.Context, id uuid.UUID) ([]coursepoint.Leaderboard, error) {
	entity := make([]coursepoint.Leaderboard, 0)

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			up.user_id,
			SUM(up.xp) AS xp,
			cs.id AS course_id
		FROM user_course_points up
		JOIN course_lessons cl ON up.lesson_id = cl.id
		JOIN course_modules cm ON cl.module_id = cm.id
		JOIN courses cs ON cm.course_id = cs.id
		WHERE cs.id = ?
		GROUP BY up.user_id, cs.id
		ORDER BY SUM(up.xp) DESC, MIN(up.updated_at) ASC
	`, id).Scan(&entity).Error

	if err != nil {
		return nil, err
	}

	return entity, nil
}
