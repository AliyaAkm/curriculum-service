package progress

import (
	"context"
	"curriculum-service/internal/domain"
	progressdomain "curriculum-service/internal/domain/progress"
	"curriculum-service/internal/repo/postgres/lessoncompletion"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type lessonInfoRow struct {
	LessonID uuid.UUID `gorm:"column:lesson_id"`
	ModuleID uuid.UUID `gorm:"column:module_id"`
	CourseID uuid.UUID `gorm:"column:course_id"`
	XPReward int       `gorm:"column:xp_reward"`
}

type subscriptionRow struct {
	UserID          uuid.UUID  `gorm:"column:user_id"`
	CourseID        uuid.UUID  `gorm:"column:course_id"`
	StartedAt       *time.Time `gorm:"column:started_at"`
	LastActivityAt  *time.Time `gorm:"column:last_activity_at"`
	CompletedAt     *time.Time `gorm:"column:completed_at"`
	CurrentLessonID *uuid.UUID `gorm:"column:current_lesson_id"`
}

type moduleProgressRow struct {
	ModuleID         uuid.UUID `gorm:"column:module_id"`
	Position         int       `gorm:"column:position"`
	TotalLessons     int       `gorm:"column:total_lessons"`
	CompletedLessons int       `gorm:"column:completed_lessons"`
}

func (r *Repo) CompleteLesson(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID) (*progressdomain.CourseProgress, error) {
	var result *progressdomain.CourseProgress

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		info, err := r.getLessonInfoTx(ctx, tx, lessonID)
		if err != nil {
			return err
		}

		hasSubscription, err := r.hasSubscriptionTx(ctx, tx, userID, info.CourseID)
		if err != nil {
			return err
		}
		if !hasSubscription {
			return domain.ErrForbidden
		}

		newlyCompleted, err := lessoncompletion.MarkTheoryAndTryComplete(ctx, tx, userID, lessonID)
		if err != nil {
			return err
		}

		result, err = r.getCourseProgressTx(ctx, tx, userID, info.CourseID)
		if result != nil {
			result.NewlyCompleted = newlyCompleted
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repo) GetCourseProgress(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*progressdomain.CourseProgress, error) {
	return r.getCourseProgressTx(ctx, r.db, userID, courseID)
}

func (r *Repo) GetLessonNotificationData(ctx context.Context, lessonID uuid.UUID) (*progressdomain.LessonNotificationData, error) {
	var row progressdomain.LessonNotificationData
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			cl.id AS lesson_id,
			COALESCE(NULLIF(lt.name, ''), 'Lesson') AS lesson_title,
			cm.course_id AS course_id,
			COALESCE(NULLIF(c.title, ''), 'Course') AS course_title
		FROM course_lessons cl
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		INNER JOIN courses c ON c.id = cm.course_id
		LEFT JOIN LATERAL (
			SELECT title.name
			FROM course_lesson_titles title
			LEFT JOIN course_locales locale ON locale.id = title.locale_id
			WHERE title.lesson_id = cl.id
			ORDER BY
				CASE locale.code
					WHEN 'ru' THEN 0
					WHEN 'en' THEN 1
					WHEN 'kk' THEN 2
					ELSE 3
				END,
				title.name
			LIMIT 1
		) lt ON TRUE
		WHERE cl.id = ?
	`, lessonID).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.LessonID == uuid.Nil {
		return nil, domain.ErrLessonNotFound
	}

	return &row, nil
}

func (r *Repo) SyncUserLevel(ctx context.Context, userID uuid.UUID) (*int, error) {
	var row struct {
		CurrentLevel  int `gorm:"column:current_level"`
		ComputedLevel int `gorm:"column:computed_level"`
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			COALESCE(level, 0)::int AS current_level,
			(
				1 + (
					COALESCE((SELECT SUM(xp) FROM user_xp_events WHERE user_id = users.id), 0)::bigint / 180
				)
			)::int AS computed_level
		FROM users
		WHERE id = ?
	`, userID).Scan(&row).Error; err != nil {
		return nil, err
	}

	if row.ComputedLevel != row.CurrentLevel {
		if err := r.db.WithContext(ctx).Exec(`
			UPDATE users
			SET level = ?
			WHERE id = ?
		`, row.ComputedLevel, userID).Error; err != nil {
			return nil, err
		}
	}

	if row.ComputedLevel > row.CurrentLevel {
		return &row.ComputedLevel, nil
	}
	return nil, nil
}

func (r *Repo) ListCourseProgress(ctx context.Context, userID uuid.UUID) ([]progressdomain.CourseProgress, error) {
	var rows []struct {
		CourseID uuid.UUID `gorm:"column:course_id"`
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT course_id
		FROM course_subscription
		WHERE user_id = ?
		  AND started_at IS NOT NULL
		ORDER BY COALESCE(last_activity_at, started_at) DESC NULLS LAST
	`, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]progressdomain.CourseProgress, 0, len(rows))
	for _, row := range rows {
		item, err := r.GetCourseProgress(ctx, userID, row.CourseID)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	return items, nil
}

func (r *Repo) getCourseProgressTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID) (*progressdomain.CourseProgress, error) {
	subscription, err := r.getSubscriptionTx(ctx, tx, userID, courseID)
	if err != nil {
		return nil, err
	}

	totalLessons, completedLessons, err := r.countCourseLessonsTx(ctx, tx, userID, courseID)
	if err != nil {
		return nil, err
	}

	modules, err := r.getModuleProgressTx(ctx, tx, userID, courseID)
	if err != nil {
		return nil, err
	}

	completedLessonIDs, err := r.getCompletedLessonIDsTx(ctx, tx, userID, courseID)
	if err != nil {
		return nil, err
	}

	theoryCompletedLessonIDs, err := r.getTheoryCompletedLessonIDsTx(ctx, tx, userID, courseID)
	if err != nil {
		return nil, err
	}

	passedQuizIDs, err := r.getPassedQuizIDsTx(ctx, tx, userID, courseID)
	if err != nil {
		return nil, err
	}

	return &progressdomain.CourseProgress{
		CourseID:                 courseID,
		UserID:                   userID,
		StartedAt:                subscription.StartedAt,
		LastActivityAt:           subscription.LastActivityAt,
		CompletedAt:              subscription.CompletedAt,
		CurrentLessonID:          subscription.CurrentLessonID,
		TotalLessons:             totalLessons,
		CompletedLessons:         completedLessons,
		ProgressPercent:          progressPercent(completedLessons, totalLessons),
		TheoryCompletedLessonIDs: theoryCompletedLessonIDs,
		CompletedLessonIDs:       completedLessonIDs,
		PassedQuizIDs:            passedQuizIDs,
		Modules:                  modules,
	}, nil
}

func (r *Repo) getLessonInfoTx(ctx context.Context, tx *gorm.DB, lessonID uuid.UUID) (lessonInfoRow, error) {
	var row lessonInfoRow
	err := tx.WithContext(ctx).Raw(`
		SELECT
			cl.id AS lesson_id,
			cl.module_id AS module_id,
			cm.course_id AS course_id,
			cl.xp_reward AS xp_reward
		FROM course_lessons cl
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		WHERE cl.id = ?
	`, lessonID).Scan(&row).Error
	if err != nil {
		return lessonInfoRow{}, err
	}
	if row.LessonID == uuid.Nil {
		return lessonInfoRow{}, domain.ErrLessonNotFound
	}

	return row, nil
}

func (r *Repo) getSubscriptionTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID) (subscriptionRow, error) {
	var row subscriptionRow
	err := tx.WithContext(ctx).Raw(`
		SELECT
			user_id,
			course_id,
			started_at,
			last_activity_at,
			completed_at,
			current_lesson_id
		FROM course_subscription
		WHERE user_id = ?
		  AND course_id = ?
	`, userID, courseID).Scan(&row).Error
	if err != nil {
		return subscriptionRow{}, err
	}
	if row.UserID == uuid.Nil {
		return subscriptionRow{}, domain.ErrCourseProgressNotFound
	}

	return row, nil
}

func (r *Repo) hasSubscriptionTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID) (bool, error) {
	var count int64
	if err := tx.WithContext(ctx).
		Table("course_subscription").
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Repo) countCourseLessonsTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID) (int, int, error) {
	var row struct {
		TotalLessons     int `gorm:"column:total_lessons"`
		CompletedLessons int `gorm:"column:completed_lessons"`
	}
	err := tx.WithContext(ctx).Raw(`
		SELECT
			COUNT(cl.id)::int AS total_lessons,
			COUNT(up.lesson_id)::int AS completed_lessons
		FROM course_modules cm
		INNER JOIN course_lessons cl ON cl.module_id = cm.id
		LEFT JOIN user_course_points up ON up.lesson_id = cl.id AND up.user_id = ?
		WHERE cm.course_id = ?
	`, userID, courseID).Scan(&row).Error
	if err != nil {
		return 0, 0, err
	}

	return row.TotalLessons, row.CompletedLessons, nil
}

func (r *Repo) getModuleProgressTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID) ([]progressdomain.ModuleProgress, error) {
	var rows []moduleProgressRow
	err := tx.WithContext(ctx).Raw(`
		SELECT
			cm.id AS module_id,
			COALESCE(cm.position, 0)::int AS position,
			COUNT(cl.id)::int AS total_lessons,
			COUNT(up.lesson_id)::int AS completed_lessons
		FROM course_modules cm
		LEFT JOIN course_lessons cl ON cl.module_id = cm.id
		LEFT JOIN user_course_points up ON up.lesson_id = cl.id AND up.user_id = ?
		WHERE cm.course_id = ?
		GROUP BY cm.id, cm.position, cm.created_at
		ORDER BY COALESCE(cm.position, 0) ASC, cm.created_at ASC
	`, userID, courseID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	items := make([]progressdomain.ModuleProgress, len(rows))
	previousModulesCompleted := true
	for i, row := range rows {
		isOpen := previousModulesCompleted
		items[i] = progressdomain.ModuleProgress{
			ModuleID:         row.ModuleID,
			Position:         row.Position,
			IsOpen:           isOpen,
			TotalLessons:     row.TotalLessons,
			CompletedLessons: row.CompletedLessons,
		}
		if row.TotalLessons > 0 && row.CompletedLessons < row.TotalLessons {
			previousModulesCompleted = false
		}
	}

	return items, nil
}

func (r *Repo) getCompletedLessonIDsTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID) ([]uuid.UUID, error) {
	var rows []struct {
		ID uuid.UUID `gorm:"column:id"`
	}
	err := tx.WithContext(ctx).Raw(`
		SELECT cl.id
		FROM user_course_points up
		INNER JOIN course_lessons cl ON cl.id = up.lesson_id
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		WHERE up.user_id = ?
		  AND cm.course_id = ?
		ORDER BY COALESCE(cm.position, 0) ASC, cl.position ASC, cl.created_at ASC
	`, userID, courseID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(rows))
	for i := range rows {
		ids[i] = rows[i].ID
	}

	return ids, nil
}

func (r *Repo) getTheoryCompletedLessonIDsTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID) ([]uuid.UUID, error) {
	var rows []struct {
		ID uuid.UUID `gorm:"column:id"`
	}
	err := tx.WithContext(ctx).Raw(`
		SELECT slp.lesson_id AS id
		FROM student_lesson_progress slp
		INNER JOIN course_lessons cl ON cl.id = slp.lesson_id
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		WHERE slp.user_id = ?
		  AND slp.theory_completed_at IS NOT NULL
		  AND cm.course_id = ?
		ORDER BY COALESCE(cm.position, 0) ASC, cl.position ASC, cl.created_at ASC
	`, userID, courseID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(rows))
	for i := range rows {
		ids[i] = rows[i].ID
	}

	return ids, nil
}

func (r *Repo) getPassedQuizIDsTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID) ([]uuid.UUID, error) {
	var rows []struct {
		ID uuid.UUID `gorm:"column:id"`
	}
	err := tx.WithContext(ctx).Raw(`
		SELECT lq.id
		FROM lesson_quiz_attempts attempt
		INNER JOIN lesson_quizzes lq ON lq.id = attempt.quiz_id
		INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		WHERE attempt.user_id = ?
		  AND cm.course_id = ?
		  AND attempt.is_correct = TRUE
		ORDER BY COALESCE(cm.position, 0) ASC, cl.position ASC, lq.position ASC
	`, userID, courseID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(rows))
	for i := range rows {
		ids[i] = rows[i].ID
	}

	return ids, nil
}

func progressPercent(completedLessons, totalLessons int) int {
	if totalLessons <= 0 {
		return 0
	}

	return int(math.Round(float64(completedLessons) / float64(totalLessons) * 100))
}
