package codeattempt

import (
	"context"
	codeattemptdomain "curriculum-service/internal/domain/codeattempt"
	practicereviewdomain "curriculum-service/internal/domain/practicereview"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repo) CreateAttempt(ctx context.Context, value codeattemptdomain.Attempt) (*codeattemptdomain.Attempt, error) {
	if value.ID == uuid.Nil {
		value.ID = uuid.New()
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Exec(`
			INSERT INTO code_execution_attempts (
				id,
				user_id,
				course_id,
				lesson_id,
				practice_id,
				run_type,
				language,
				passed,
				error_type,
				error_message,
				output,
				duration_ms,
				code_hash,
				xp_awarded
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
		`, value.ID,
			value.UserID,
			value.CourseID,
			value.LessonID,
			value.PracticeID,
			value.RunType,
			value.Language,
			value.Passed,
			value.ErrorType,
			value.ErrorMessage,
			value.Output,
			value.DurationMS,
			value.CodeHash,
		).Error; err != nil {
			return err
		}

		practiceID, err := uuid.Parse(value.PracticeID)
		if err != nil {
			return err
		}

		progressStatus := practicereviewdomain.ProgressStatusInProgress
		if value.RunType == codeattemptdomain.RunTypeSubmit && value.Passed {
			progressStatus = practicereviewdomain.ProgressStatusCompleted
		}

		if value.CourseID != nil && value.LessonID != nil {
			if err := upsertPracticeProgressByAttemptTx(ctx, tx, value, practiceID, progressStatus); err != nil {
				return err
			}
		}

		if value.RunType != codeattemptdomain.RunTypeSubmit || !value.Passed || value.XPReward <= 0 {
			return updateCourseActivityByPracticeTx(ctx, tx, value)
		}

		result := tx.WithContext(ctx).Exec(`
			INSERT INTO practice_xp_awards (
				id,
				user_id,
				practice_id,
				course_id,
				lesson_id,
				xp
			)
			VALUES (?, ?, ?, ?, ?, ?)
			ON CONFLICT (user_id, practice_id) DO NOTHING
		`, uuid.New(), value.UserID, practiceID, value.CourseID, value.LessonID, value.XPReward)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			value.XPAwarded = value.XPReward
			if err := tx.WithContext(ctx).Exec(`
				UPDATE code_execution_attempts
				SET xp_awarded = ?
				WHERE id = ?
			`, value.XPAwarded, value.ID).Error; err != nil {
				return err
			}
			if err := syncUserLevelByPracticeTx(ctx, tx, value.UserID); err != nil {
				return err
			}
		}

		return updateCourseActivityByPracticeTx(ctx, tx, value)
	})
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func upsertPracticeProgressByAttemptTx(ctx context.Context, tx *gorm.DB, value codeattemptdomain.Attempt, practiceID uuid.UUID, status string) error {
	completedExpr := "student_practice_progress.completed_at"
	if status == practicereviewdomain.ProgressStatusCompleted {
		completedExpr = "COALESCE(student_practice_progress.completed_at, NOW())"
	}

	return tx.WithContext(ctx).Exec(`
		INSERT INTO student_practice_progress (
			id,
			student_id,
			practice_id,
			course_id,
			lesson_id,
			status,
			started_at,
			completed_at,
			last_attempt_at,
			attempts_count
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			NOW(),
			CASE WHEN ? = 'completed' THEN NOW() ELSE NULL END,
			NOW(),
			1
		)
		ON CONFLICT (student_id, practice_id) DO UPDATE
		SET status = CASE
				WHEN student_practice_progress.status = 'completed' THEN student_practice_progress.status
				ELSE EXCLUDED.status
			END,
		    completed_at = `+completedExpr+`,
		    last_attempt_at = NOW(),
		    attempts_count = student_practice_progress.attempts_count + 1,
		    updated_at = NOW()
	`, uuid.New(),
		value.UserID,
		practiceID,
		*value.CourseID,
		*value.LessonID,
		status,
		status,
	).Error
}

func updateCourseActivityByPracticeTx(ctx context.Context, tx *gorm.DB, value codeattemptdomain.Attempt) error {
	if value.CourseID == nil || value.LessonID == nil {
		return nil
	}

	return tx.WithContext(ctx).Exec(`
		UPDATE course_subscription
		SET started_at = COALESCE(started_at, NOW()),
		    last_activity_at = NOW(),
		    current_lesson_id = ?
		WHERE user_id = ?
		  AND course_id = ?
	`, *value.LessonID, value.UserID, *value.CourseID).Error
}

func syncUserLevelByPracticeTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
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
