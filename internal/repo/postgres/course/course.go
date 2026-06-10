package course

import (
	"context"
	"curriculum-service/internal/domain"
	"curriculum-service/internal/domain/course"
	"curriculum-service/internal/domain/price"
	"curriculum-service/internal/domain/reviewlog"
	dtocourse "curriculum-service/internal/http/dto/course"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"

	"strings"
)

func (r *Repo) GetAllCourses(ctx context.Context, query dtocourse.GetCoursesQuery) ([]course.Course, error) {
	var courses []course.Course
	db := r.db.WithContext(ctx).
		Preload("Status").
		Preload("Level").
		Preload("DurationCategory").
		Preload("Author.Roles").
		Preload("Tags").
		Preload("Topic")

	db = db.Joins("LEFT JOIN course_topics ct ON ct.id=courses.topic_id").
		Joins("LEFT JOIN course_levels cl ON cl.id = courses.level_id").
		Joins("LEFT JOIN course_statuses cs ON cs.id = courses.status_id").
		Joins("LEFT JOIN course_duration_categories cdc ON cdc.id = courses.duration_category_id").
		Joins("LEFT JOIN users u ON u.id = courses.author_id").
		Joins("LEFT JOIN course_course_tags cct ON cct.course_id = courses.id").
		Joins("LEFT JOIN course_tags tg ON tg.id = cct.tag_id")

	if strings.TrimSpace(query.Search) != "" {
		search := "%" + strings.TrimSpace(query.Search) + "%"
		// to do: поменять email на логин в фильтрации
		db = db.Where(`
			courses.title ILIKE ?
			OR courses.subtitle ILIKE ?
			OR courses.description ILIKE ?
			OR ct.name ILIKE ?
			OR ct.code ILIKE ?
			OR cl.name ILIKE ?
			OR cl.code ILIKE ?
			OR cdc.name ILIKE ?
			OR cdc.code ILIKE ?
			OR tg.name ILIKE ?
			OR tg.code ILIKE ?
			OR u.email ILIKE ? 
			OR cs.name ILIKE ?
			OR cs.code ILIKE ?
		`,
			search,
			search,
			search,
			search,
			search,
			search,
			search,
			search,
			search,
			search,
			search,
			search,
			search,
			search,
		)
	}
	if query.Level != "" {
		db = db.Where("cl.code = ?", query.Level)
	}
	if query.DurationCategory != "" {
		db = db.Where("cdc.code = ?", query.DurationCategory)
	}

	if query.MinRating != nil {
		db = db.Where("courses.rating >= ?", *query.MinRating)
	}

	if query.HasCertificate != nil {
		db = db.Where("courses.has_certificate = ?", *query.HasCertificate)
	}
	if query.Topic != "" {
		db = db.Where("ct.code = ?", query.Topic)
	}

	db = db.Distinct("courses.*")

	limit := 10
	if query.Limit > 0 {
		limit = query.Limit
	}
	page := 1
	if query.Page > 0 {
		page = query.Page
	}
	offset := (page - 1) * limit

	err := db.Limit(limit).Offset(offset).Find(&courses).Error
	if err != nil {
		return nil, err
	}

	return courses, nil
}
func (r *Repo) CreateCourse(ctx context.Context, value *course.Course) (uuid.UUID, error) {
	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return uuid.UUID{}, err
	}

	return value.ID, nil
}

func (r *Repo) GetPendingCheckCourses(ctx context.Context) ([]course.Course, error) {
	var courses []course.Course

	err := r.db.WithContext(ctx).
		Preload("Level").
		Preload("Topic").
		Where("courses.is_checked = ?", false).
		Order("courses.created_at DESC").
		Find(&courses).Error
	if err != nil {
		return nil, err
	}

	return courses, nil
}

// создать подписку пользователю при оплате
func (r *Repo) CreateSubscription(ctx context.Context, value *course.Subscription) error {
	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return err
	}
	return nil
}

// получить подписку по айди подписке
func (r *Repo) GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*course.Subscription, error) {
	var entity course.Subscription
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// проверить подписку
func (r *Repo) HasSubscription(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&course.Subscription{}).Where("user_id = ? AND course_id = ?", userID, courseID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repo) ReviewCourse(ctx context.Context, log *reviewlog.CourseReviewLog, isApproved bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// existence check -> 404
		var existing course.Course
		err := tx.Select("id").Where("id = ?", log.CourseID).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrCourseNotFound
		}
		if err != nil {
			return err
		}

		// NOTE: map update (not struct) is required so is_approved=false actually persists.
		// A struct update would skip the zero value `false` and silently leave is_approved unchanged.
		if err := tx.Model(&course.Course{}).
			Where("id = ?", log.CourseID).
			Updates(map[string]any{
				"is_checked":  true,
				"is_approved": isApproved,
				"updated_at":  time.Now(),
			}).Error; err != nil {
			return err
		}

		return tx.Create(log).Error
	})
}

func (r *Repo) SetCourseUnchecked(ctx context.Context, courseID uuid.UUID) error {
	// map update (not struct): a struct update would skip the zero values `false`
	// and leave is_checked / is_approved unchanged.
	return r.db.WithContext(ctx).
		Model(&course.Course{}).
		Where("id = ?", courseID).
		Updates(map[string]any{
			"is_checked":  false,
			"is_approved": false,
			"updated_at":  time.Now(),
		}).Error
}

func (r *Repo) GetLatestReviewLog(ctx context.Context, courseID uuid.UUID) (*reviewlog.CourseReviewLog, error) {
	var log reviewlog.CourseReviewLog
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("created_at DESC").
		First(&log).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // no review history yet -> not an error
	}
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *Repo) UpdateCourse(ctx context.Context, id uuid.UUID, value *course.Course) error {
	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(value).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) GetCourseByID(ctx context.Context, id uuid.UUID) (*course.Course, error) {
	var entity course.Course
	err := r.courseQuery(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repo) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	var entity course.Course
	err := r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteCoursePrice(ctx context.Context, courseID uuid.UUID) error {
	var entity price.CoursePrice
	err := r.db.WithContext(ctx).Delete(&entity, "course_id = ?", courseID).Error
	if err != nil {
		return err
	}
	return nil

}

func (r *Repo) courseQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Preload("Status").Preload("DurationCategory").Preload("Level").Preload("Author.Roles").Preload("Tags").Preload("Topic")
}
