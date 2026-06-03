package certificate

import (
	"context"
	"curriculum-service/internal/domain"
	certificatedomain "curriculum-service/internal/domain/certificate"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repo) FindByUserAndCourse(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*certificatedomain.Certificate, error) {
	var value certificatedomain.Certificate
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&value).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrCertificateNotFound
		}
		return nil, err
	}

	return &value, nil
}

func (r *Repo) GetCourseCompletion(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (certificatedomain.CourseCompletion, error) {
	var row certificatedomain.CourseCompletion
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			c.id AS course_id,
			c.title AS course_title,
			COALESCE(u.login, '') AS user_login,
			c.has_certificate AS has_certificate,
			cs.completed_at AS completed_at,
			COUNT(cl.id)::int AS total_lessons,
			COUNT(up.lesson_id)::int AS completed_lessons
		FROM courses c
		LEFT JOIN users u ON u.id = ?
		LEFT JOIN course_subscription cs ON cs.course_id = c.id AND cs.user_id = ?
		LEFT JOIN course_modules cm ON cm.course_id = c.id
		LEFT JOIN course_lessons cl ON cl.module_id = cm.id
		LEFT JOIN user_course_points up ON up.lesson_id = cl.id AND up.user_id = ?
		WHERE c.id = ?
		GROUP BY c.id, c.title, u.login, c.has_certificate, cs.completed_at
	`, userID, userID, userID, courseID).Scan(&row).Error
	if err != nil {
		return certificatedomain.CourseCompletion{}, err
	}
	if row.CourseID == uuid.Nil {
		return certificatedomain.CourseCompletion{}, domain.ErrCourseNotFound
	}

	return row, nil
}

func (r *Repo) Create(ctx context.Context, value *certificatedomain.Certificate) error {
	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return err
	}

	return nil
}
