package lessoncompletion

import (
	"context"
	"curriculum-service/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type lessonInfo struct {
	LessonID uuid.UUID `gorm:"column:lesson_id"`
	CourseID uuid.UUID `gorm:"column:course_id"`
	XPReward int       `gorm:"column:xp_reward"`
}

func MarkTheoryAndTryComplete(ctx context.Context, tx *gorm.DB, userID uuid.UUID, lessonID uuid.UUID) (bool, error) {
	info, err := getLessonInfo(ctx, tx, lessonID)
	if err != nil {
		return false, err
	}

	hasSubscription, err := hasSubscription(ctx, tx, userID, info.CourseID)
	if err != nil {
		return false, err
	}
	if !hasSubscription {
		return false, domain.ErrForbidden
	}

	if err := tx.WithContext(ctx).Exec(`
		INSERT INTO student_lesson_progress (
			id,
			user_id,
			lesson_id,
			course_id,
			theory_completed_at,
			last_activity_at
		)
		VALUES (?, ?, ?, ?, NOW(), NOW())
		ON CONFLICT (user_id, lesson_id) DO UPDATE
		SET theory_completed_at = COALESCE(student_lesson_progress.theory_completed_at, NOW()),
		    last_activity_at = NOW(),
		    updated_at = NOW()
	`, uuid.New(), userID, info.LessonID, info.CourseID).Error; err != nil {
		return false, err
	}

	if err := touchCourseActivity(ctx, tx, userID, info.CourseID, info.LessonID); err != nil {
		return false, err
	}

	return TryComplete(ctx, tx, userID, info.LessonID)
}

func TryCompleteByQuiz(ctx context.Context, tx *gorm.DB, userID uuid.UUID, quizID uuid.UUID) (bool, error) {
	var row struct {
		LessonID uuid.UUID `gorm:"column:lesson_id"`
	}
	if err := tx.WithContext(ctx).Raw(`
		SELECT lesson_id
		FROM lesson_quizzes
		WHERE id = ?
	`, quizID).Scan(&row).Error; err != nil {
		return false, err
	}
	if row.LessonID == uuid.Nil {
		return false, domain.ErrQuizNotFound
	}
	return TryComplete(ctx, tx, userID, row.LessonID)
}

func TryCompleteByPractice(ctx context.Context, tx *gorm.DB, userID uuid.UUID, practiceID uuid.UUID) (bool, error) {
	var row struct {
		LessonID uuid.UUID `gorm:"column:lesson_id"`
	}
	if err := tx.WithContext(ctx).Raw(`
		SELECT lesson_id
		FROM practice_tasks
		WHERE id = ?
	`, practiceID).Scan(&row).Error; err != nil {
		return false, err
	}
	if row.LessonID == uuid.Nil {
		return false, domain.ErrPracticeNotFound
	}
	return TryComplete(ctx, tx, userID, row.LessonID)
}

func CanStartPractice(ctx context.Context, tx *gorm.DB, userID uuid.UUID, practiceID uuid.UUID) (bool, error) {
	var row struct {
		Ready bool `gorm:"column:ready"`
	}
	if err := tx.WithContext(ctx).Raw(`
		WITH info AS (
			SELECT
				pt.id AS practice_id,
				pt.lesson_id,
				cm.course_id
			FROM practice_tasks pt
			INNER JOIN course_lessons cl ON cl.id = pt.lesson_id
			INNER JOIN course_modules cm ON cm.id = cl.module_id
			WHERE pt.id = ?
		)
		SELECT EXISTS (
			SELECT 1
			FROM info
			INNER JOIN course_subscription cs
				ON cs.course_id = info.course_id
			   AND cs.user_id = ?
			WHERE (
				EXISTS (
					SELECT 1
					FROM user_course_points ucp
					WHERE ucp.user_id = ?
					  AND ucp.lesson_id = info.lesson_id
				)
				OR (
					EXISTS (
						SELECT 1
						FROM student_lesson_progress slp
						WHERE slp.user_id = ?
						  AND slp.lesson_id = info.lesson_id
						  AND slp.theory_completed_at IS NOT NULL
					)
					AND NOT EXISTS (
						SELECT 1
						FROM lesson_quizzes lq
						WHERE lq.lesson_id = info.lesson_id
						  AND NOT EXISTS (
							SELECT 1
							FROM lesson_quiz_attempts lqa
							WHERE lqa.user_id = ?
							  AND lqa.quiz_id = lq.id
							  AND lqa.is_correct = TRUE
						  )
					)
				)
			)
		) AS ready
	`, practiceID, userID, userID, userID, userID).Scan(&row).Error; err != nil {
		return false, err
	}
	return row.Ready, nil
}

func TryComplete(ctx context.Context, tx *gorm.DB, userID uuid.UUID, lessonID uuid.UUID) (bool, error) {
	info, err := getLessonInfo(ctx, tx, lessonID)
	if err != nil {
		return false, err
	}

	var ready struct {
		Ready bool `gorm:"column:ready"`
	}
	if err := tx.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM student_lesson_progress slp
			WHERE slp.user_id = ?
			  AND slp.lesson_id = ?
			  AND slp.theory_completed_at IS NOT NULL
			  AND NOT EXISTS (
				SELECT 1
				FROM lesson_quizzes lq
				WHERE lq.lesson_id = slp.lesson_id
				  AND NOT EXISTS (
					SELECT 1
					FROM lesson_quiz_attempts lqa
					WHERE lqa.user_id = slp.user_id
					  AND lqa.quiz_id = lq.id
					  AND lqa.is_correct = TRUE
				  )
			  )
			  AND NOT EXISTS (
				SELECT 1
				FROM practice_tasks pt
				WHERE pt.lesson_id = slp.lesson_id
				  AND NOT EXISTS (
					SELECT 1
					FROM student_practice_progress spp
					WHERE spp.student_id = slp.user_id
					  AND spp.practice_id = pt.id
					  AND spp.status = 'completed'
				  )
			  )
		) AS ready
	`, userID, info.LessonID).Scan(&ready).Error; err != nil {
		return false, err
	}
	if !ready.Ready {
		return false, nil
	}

	if err := tx.WithContext(ctx).Exec(`
		UPDATE student_lesson_progress
		SET completed_at = COALESCE(completed_at, NOW()),
		    last_activity_at = NOW(),
		    updated_at = NOW()
		WHERE user_id = ?
		  AND lesson_id = ?
	`, userID, info.LessonID).Error; err != nil {
		return false, err
	}

	insert := tx.WithContext(ctx).Exec(`
		INSERT INTO user_course_points (
			id,
			lesson_id,
			user_id,
			xp
		)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (user_id, lesson_id) DO NOTHING
	`, uuid.New(), info.LessonID, userID, info.XPReward)
	if insert.Error != nil {
		return false, insert.Error
	}

	if err := syncCourseCompletion(ctx, tx, userID, info.CourseID, info.LessonID); err != nil {
		return false, err
	}
	if insert.RowsAffected > 0 {
		if err := syncUserLevel(ctx, tx, userID); err != nil {
			return false, err
		}
	}

	return insert.RowsAffected > 0, nil
}

func getLessonInfo(ctx context.Context, tx *gorm.DB, lessonID uuid.UUID) (lessonInfo, error) {
	var row lessonInfo
	if err := tx.WithContext(ctx).Raw(`
		SELECT
			cl.id AS lesson_id,
			cm.course_id AS course_id,
			cl.xp_reward AS xp_reward
		FROM course_lessons cl
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		WHERE cl.id = ?
	`, lessonID).Scan(&row).Error; err != nil {
		return lessonInfo{}, err
	}
	if row.LessonID == uuid.Nil {
		return lessonInfo{}, domain.ErrLessonNotFound
	}
	return row, nil
}

func hasSubscription(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID) (bool, error) {
	var row struct {
		Exists bool `gorm:"column:exists"`
	}
	if err := tx.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM course_subscription
			WHERE user_id = ?
			  AND course_id = ?
		) AS exists
	`, userID, courseID).Scan(&row).Error; err != nil {
		return false, err
	}
	return row.Exists, nil
}

func touchCourseActivity(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID, lessonID uuid.UUID) error {
	return tx.WithContext(ctx).Exec(`
		UPDATE course_subscription
		SET started_at = COALESCE(started_at, NOW()),
		    last_activity_at = NOW(),
		    current_lesson_id = ?
		WHERE user_id = ?
		  AND course_id = ?
	`, lessonID, userID, courseID).Error
}

func syncCourseCompletion(ctx context.Context, tx *gorm.DB, userID uuid.UUID, courseID uuid.UUID, lessonID uuid.UUID) error {
	return tx.WithContext(ctx).Exec(`
		UPDATE course_subscription
		SET started_at = COALESCE(started_at, NOW()),
		    last_activity_at = NOW(),
		    current_lesson_id = ?,
		    completed_at = CASE
				WHEN (
					SELECT COUNT(*)::int
					FROM course_modules cm
					INNER JOIN course_lessons cl ON cl.module_id = cm.id
					WHERE cm.course_id = ?
				) > 0
				AND (
					SELECT COUNT(DISTINCT ucp.lesson_id)::int
					FROM user_course_points ucp
					INNER JOIN course_lessons cl ON cl.id = ucp.lesson_id
					INNER JOIN course_modules cm ON cm.id = cl.module_id
					WHERE ucp.user_id = ?
					  AND cm.course_id = ?
				) >= (
					SELECT COUNT(*)::int
					FROM course_modules cm
					INNER JOIN course_lessons cl ON cl.module_id = cm.id
					WHERE cm.course_id = ?
				)
				THEN COALESCE(completed_at, NOW())
				ELSE completed_at
			END
		WHERE user_id = ?
		  AND course_id = ?
	`, lessonID, courseID, userID, courseID, courseID, userID, courseID).Error
}

func syncUserLevel(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
	return tx.WithContext(ctx).Exec(`
		UPDATE users
		SET level = GREATEST(
			COALESCE(level, 1),
			1 + (
				COALESCE((SELECT SUM(xp) FROM user_xp_events WHERE user_id = users.id), 0)::bigint / 180
			)::int
		)
		WHERE id = ?
	`, userID).Error
}
