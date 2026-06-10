package practicereview

import (
	"context"
	"curriculum-service/internal/domain"
	practicedomain "curriculum-service/internal/domain/practice"
	practicereviewdomain "curriculum-service/internal/domain/practicereview"
	"curriculum-service/internal/repo/postgres/lessoncompletion"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type practiceInfoRow struct {
	PracticeID uuid.UUID `gorm:"column:practice_id"`
	CourseID   uuid.UUID `gorm:"column:course_id"`
	LessonID   uuid.UUID `gorm:"column:lesson_id"`
	Language   string    `gorm:"column:language"`
	XPReward   int       `gorm:"column:xp_reward"`
	CheckType  string    `gorm:"column:check_type"`
}

func (r *Repo) CreateSubmission(ctx context.Context, req practicereviewdomain.CreateSubmissionRequest) (*practicereviewdomain.Submission, error) {
	var createdID uuid.UUID
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		info, err := getPracticeInfoTx(ctx, tx, req.PracticeID)
		if err != nil {
			return err
		}
		if info.CheckType != practicedomain.CheckTypeManual {
			return domain.ErrPracticeManualReviewNotAllowed
		}

		hasSubscription, err := hasSubscriptionTx(ctx, tx, req.StudentID, info.CourseID)
		if err != nil {
			return err
		}
		if !hasSubscription {
			return domain.ErrCourseSubscriptionNotFound
		}
		canStart, err := lessoncompletion.CanStartPractice(ctx, tx, req.StudentID, req.PracticeID)
		if err != nil {
			return err
		}
		if !canStart {
			return domain.ErrPracticePrerequisitesNotMet
		}

		completed, err := isPracticeCompletedTx(ctx, tx, req.StudentID, req.PracticeID)
		if err != nil {
			return err
		}
		if completed {
			return domain.ErrPracticeAlreadyCompleted
		}

		hasActive, err := hasActiveSubmissionTx(ctx, tx, req.StudentID, req.PracticeID)
		if err != nil {
			return err
		}
		if hasActive {
			return domain.ErrPracticeSubmissionExists
		}

		attemptNumber, err := nextAttemptNumberTx(ctx, tx, req.StudentID, req.PracticeID)
		if err != nil {
			return err
		}

		language := strings.TrimSpace(req.Language)
		if language == "" {
			language = info.Language
		}

		createdID = uuid.New()
		if err = tx.WithContext(ctx).Exec(`
			INSERT INTO practice_review_submissions (
				id,
				practice_id,
				student_id,
				course_id,
				lesson_id,
				status,
				code,
				language,
				output,
				error,
				error_type,
				exit_code,
				duration_ms,
				attempt_number
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, createdID,
			req.PracticeID,
			req.StudentID,
			info.CourseID,
			info.LessonID,
			practicereviewdomain.SubmissionStatusSubmitted,
			req.Code,
			language,
			req.Output,
			req.Error,
			req.ErrorType,
			req.ExitCode,
			req.DurationMS,
			attemptNumber,
		).Error; err != nil {
			return err
		}

		return upsertPracticeProgressTx(ctx, tx, progressUpdate{
			StudentID:     req.StudentID,
			PracticeID:    req.PracticeID,
			CourseID:      info.CourseID,
			LessonID:      info.LessonID,
			Status:        practicereviewdomain.ProgressStatusSubmitted,
			IncrementRuns: true,
		})
	})
	if err != nil {
		return nil, err
	}

	return r.getSubmission(ctx, r.db, createdID)
}

func (r *Repo) ListStudentSubmissions(ctx context.Context, filter practicereviewdomain.StudentListFilter) ([]practicereviewdomain.Submission, error) {
	args := []any{filter.StudentID}
	query := baseSubmissionSelect() + `
		WHERE prs.student_id = ?
	`
	query, args = appendSubmissionFilters(query, args, filter.CourseID, filter.PracticeID, nil, filter.Status)
	query += `
		ORDER BY prs.created_at DESC
	`

	var rows []practicereviewdomain.Submission
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	if rows == nil {
		rows = []practicereviewdomain.Submission{}
	}
	return rows, err
}

func (r *Repo) GetStudentSubmission(ctx context.Context, studentID uuid.UUID, submissionID uuid.UUID) (*practicereviewdomain.Submission, error) {
	return r.getSubmissionWithAccess(ctx, r.db, submissionID, `
		prs.student_id = ?
	`, studentID)
}

func (r *Repo) ListTeacherSubmissions(ctx context.Context, filter practicereviewdomain.TeacherListFilter) ([]practicereviewdomain.Submission, error) {
	args := []any{}
	query := baseSubmissionSelect() + `
		WHERE 1 = 1
	`
	if !filter.IsAdmin {
		query += `
			AND c.author_id = ?
		`
		args = append(args, filter.TeacherID)
	}
	query, args = appendSubmissionFilters(query, args, filter.CourseID, filter.PracticeID, filter.StudentID, filter.Status)
	query += `
		ORDER BY
			CASE prs.status
				WHEN 'submitted' THEN 0
				WHEN 'in_review' THEN 1
				WHEN 'changes_requested' THEN 2
				ELSE 3
			END,
			prs.created_at ASC
		LIMIT ? OFFSET ?
	`
	args = append(args, filter.Limit, filter.Offset)

	var rows []practicereviewdomain.Submission
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	if rows == nil {
		rows = []practicereviewdomain.Submission{}
	}
	return rows, err
}

func (r *Repo) GetTeacherSubmission(ctx context.Context, teacherID uuid.UUID, isAdmin bool, submissionID uuid.UUID) (*practicereviewdomain.Submission, error) {
	var result *practicereviewdomain.Submission
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var submission *practicereviewdomain.Submission
		var err error
		if isAdmin {
			submission, err = r.getSubmission(ctx, tx, submissionID)
		} else {
			submission, err = r.getSubmissionWithAccess(ctx, tx, submissionID, "c.author_id = ?", teacherID)
		}
		if err != nil {
			return err
		}

		if submission.Status == practicereviewdomain.SubmissionStatusSubmitted {
			if err = tx.WithContext(ctx).Exec(`
				UPDATE practice_review_submissions
				SET status = ?,
				    updated_at = NOW()
				WHERE id = ?
				  AND status = ?
			`, practicereviewdomain.SubmissionStatusInReview, submissionID, practicereviewdomain.SubmissionStatusSubmitted).Error; err != nil {
				return err
			}
			submission.Status = practicereviewdomain.SubmissionStatusInReview
		}
		result = submission
		return nil
	})
	return result, err
}

func (r *Repo) ReviewSubmission(ctx context.Context, req practicereviewdomain.ReviewSubmissionRequest) (*practicereviewdomain.Submission, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		submission, err := r.getSubmissionForReviewTx(ctx, tx, req.SubmissionID, req.TeacherID, req.IsAdmin)
		if err != nil {
			return err
		}

		if submission.Status == practicereviewdomain.SubmissionStatusApproved {
			return domain.ErrPracticeAlreadyCompleted
		}

		if err = tx.WithContext(ctx).Exec(`
			UPDATE practice_review_submissions
			SET status = ?,
			    teacher_comment = ?,
			    reviewed_by = ?,
			    reviewed_at = NOW(),
			    updated_at = NOW()
			WHERE id = ?
		`, req.Status, req.Comment, req.TeacherID, req.SubmissionID).Error; err != nil {
			return err
		}

		progressStatus := practicereviewdomain.ProgressStatusChangesRequested
		if req.Status == practicereviewdomain.SubmissionStatusApproved {
			progressStatus = practicereviewdomain.ProgressStatusCompleted
		}

		if err = upsertPracticeProgressTx(ctx, tx, progressUpdate{
			StudentID:            submission.StudentID,
			PracticeID:           submission.PracticeID,
			CourseID:             submission.CourseID,
			LessonID:             submission.LessonID,
			Status:               progressStatus,
			ApprovedSubmissionID: req.SubmissionID,
		}); err != nil {
			return err
		}

		if req.Status != practicereviewdomain.SubmissionStatusApproved {
			return nil
		}

		info, err := getPracticeInfoTx(ctx, tx, submission.PracticeID)
		if err != nil {
			return err
		}

		if info.XPReward > 0 {
			result := tx.WithContext(ctx).Exec(`
				INSERT INTO practice_xp_awards (
					id,
					user_id,
					practice_id,
					course_id,
					lesson_id,
					xp,
					submission_id
				)
				VALUES (?, ?, ?, ?, ?, ?, ?)
				ON CONFLICT (user_id, practice_id) DO NOTHING
			`, uuid.New(), submission.StudentID, submission.PracticeID, submission.CourseID, submission.LessonID, info.XPReward, req.SubmissionID)
			if result.Error != nil {
				return result.Error
			}
		}

		if _, err = lessoncompletion.TryCompleteByPractice(ctx, tx, submission.StudentID, submission.PracticeID); err != nil {
			return err
		}

		return syncUserLevelTx(ctx, tx, submission.StudentID)
	})
	if err != nil {
		return nil, err
	}

	return r.getSubmission(ctx, r.db, req.SubmissionID)
}

func (r *Repo) getSubmission(ctx context.Context, tx *gorm.DB, submissionID uuid.UUID) (*practicereviewdomain.Submission, error) {
	return r.getSubmissionWithAccess(ctx, tx, submissionID, "1 = 1")
}

func (r *Repo) getSubmissionWithAccess(ctx context.Context, tx *gorm.DB, submissionID uuid.UUID, accessWhere string, accessArgs ...any) (*practicereviewdomain.Submission, error) {
	args := []any{submissionID}
	args = append(args, accessArgs...)
	var row practicereviewdomain.Submission
	err := tx.WithContext(ctx).Raw(baseSubmissionSelect()+`
		WHERE prs.id = ?
		  AND `+accessWhere+`
	`, args...).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == uuid.Nil {
		return nil, domain.ErrPracticeSubmissionNotFound
	}
	return &row, nil
}

func (r *Repo) getSubmissionForReviewTx(ctx context.Context, tx *gorm.DB, submissionID uuid.UUID, teacherID uuid.UUID, isAdmin bool) (*practicereviewdomain.Submission, error) {
	if isAdmin {
		return r.getSubmission(ctx, tx, submissionID)
	}
	return r.getSubmissionWithAccess(ctx, tx, submissionID, "c.author_id = ?", teacherID)
}

func appendSubmissionFilters(query string, args []any, courseID *uuid.UUID, practiceID *uuid.UUID, studentID *uuid.UUID, status string) (string, []any) {
	if courseID != nil {
		query += `
			AND prs.course_id = ?
		`
		args = append(args, *courseID)
	}
	if practiceID != nil {
		query += `
			AND prs.practice_id = ?
		`
		args = append(args, *practiceID)
	}
	if studentID != nil {
		query += `
			AND prs.student_id = ?
		`
		args = append(args, *studentID)
	}
	if status != "" {
		query += `
			AND prs.status = ?
		`
		args = append(args, status)
	}
	return query, args
}

func baseSubmissionSelect() string {
	return `
		SELECT
			prs.id,
			prs.practice_id,
			prs.student_id,
			prs.course_id,
			prs.lesson_id,
			prs.status,
			prs.code,
			prs.language,
			prs.output,
			prs.error,
			prs.error_type,
			prs.exit_code,
			prs.duration_ms,
			prs.teacher_comment,
			prs.reviewed_by,
			prs.reviewed_at,
			prs.attempt_number,
			pt.title AS practice_title,
			COALESCE(u.email, '') AS student_email,
			c.title AS course_title,
			COALESCE(lt.name, '') AS lesson_title,
			spp.status AS progress_status,
			spp.started_at AS progress_started_at,
			spp.completed_at AS progress_completed_at,
			spp.last_attempt_at AS progress_last_attempt_at,
			COALESCE(spp.attempts_count, 0)::int AS progress_attempts_count,
			prs.created_at,
			prs.updated_at
		FROM practice_review_submissions prs
		INNER JOIN practice_tasks pt ON pt.id = prs.practice_id
		INNER JOIN course_lessons cl ON cl.id = prs.lesson_id
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		INNER JOIN courses c ON c.id = prs.course_id
		LEFT JOIN users u ON u.id = prs.student_id
		LEFT JOIN student_practice_progress spp ON spp.student_id = prs.student_id AND spp.practice_id = prs.practice_id
		LEFT JOIN LATERAL (
			SELECT title.name
			FROM course_lesson_titles title
			WHERE title.lesson_id = cl.id
			ORDER BY title.name
			LIMIT 1
		) lt ON TRUE
	`
}

func getPracticeInfoTx(ctx context.Context, tx *gorm.DB, practiceID uuid.UUID) (practiceInfoRow, error) {
	var row practiceInfoRow
	err := tx.WithContext(ctx).Raw(`
		SELECT
			pt.id AS practice_id,
			cm.course_id AS course_id,
			pt.lesson_id AS lesson_id,
			pt.language AS language,
			pt.xp_reward AS xp_reward,
			pt.check_type AS check_type
		FROM practice_tasks pt
		INNER JOIN course_lessons cl ON cl.id = pt.lesson_id
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		WHERE pt.id = ?
	`, practiceID).Scan(&row).Error
	if err != nil {
		return practiceInfoRow{}, err
	}
	if row.PracticeID == uuid.Nil {
		return practiceInfoRow{}, domain.ErrPracticeNotFound
	}
	return row, nil
}

func hasSubscriptionTx(ctx context.Context, tx *gorm.DB, studentID uuid.UUID, courseID uuid.UUID) (bool, error) {
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
	`, studentID, courseID).Scan(&row).Error; err != nil {
		return false, err
	}
	return row.Exists, nil
}

func isPracticeCompletedTx(ctx context.Context, tx *gorm.DB, studentID uuid.UUID, practiceID uuid.UUID) (bool, error) {
	var row struct {
		Completed bool `gorm:"column:completed"`
	}
	if err := tx.WithContext(ctx).Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM student_practice_progress
			WHERE student_id = ?
			  AND practice_id = ?
			  AND status = ?
		) OR EXISTS (
			SELECT 1
			FROM practice_xp_awards
			WHERE user_id = ?
			  AND practice_id = ?
		) AS completed
	`, studentID, practiceID, practicereviewdomain.ProgressStatusCompleted, studentID, practiceID).Scan(&row).Error; err != nil {
		return false, err
	}
	return row.Completed, nil
}

func hasActiveSubmissionTx(ctx context.Context, tx *gorm.DB, studentID uuid.UUID, practiceID uuid.UUID) (bool, error) {
	var count int64
	if err := tx.WithContext(ctx).
		Table("practice_review_submissions").
		Where("student_id = ? AND practice_id = ? AND status IN ?", studentID, practiceID, []string{
			practicereviewdomain.SubmissionStatusSubmitted,
			practicereviewdomain.SubmissionStatusInReview,
		}).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func nextAttemptNumberTx(ctx context.Context, tx *gorm.DB, studentID uuid.UUID, practiceID uuid.UUID) (int, error) {
	var row struct {
		NextAttempt int `gorm:"column:next_attempt"`
	}
	err := tx.WithContext(ctx).Raw(`
		SELECT COALESCE(MAX(attempt_number), 0)::int + 1 AS next_attempt
		FROM practice_review_submissions
		WHERE student_id = ?
		  AND practice_id = ?
	`, studentID, practiceID).Scan(&row).Error
	if err != nil {
		return 0, err
	}
	if row.NextAttempt <= 0 {
		row.NextAttempt = 1
	}
	return row.NextAttempt, nil
}

type progressUpdate struct {
	StudentID            uuid.UUID
	PracticeID           uuid.UUID
	CourseID             uuid.UUID
	LessonID             uuid.UUID
	Status               string
	IncrementRuns        bool
	ApprovedSubmissionID uuid.UUID
}

func upsertPracticeProgressTx(ctx context.Context, tx *gorm.DB, update progressUpdate) error {
	completedExpr := "student_practice_progress.completed_at"
	if update.Status == practicereviewdomain.ProgressStatusCompleted {
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
			attempts_count,
			approved_submission_id
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
			CASE WHEN ? THEN 1 ELSE 0 END,
			NULLIF(?, '00000000-0000-0000-0000-000000000000')::uuid
		)
		ON CONFLICT (student_id, practice_id) DO UPDATE
		SET status = CASE
				WHEN student_practice_progress.status = 'completed' THEN student_practice_progress.status
				ELSE EXCLUDED.status
			END,
		    completed_at = `+completedExpr+`,
		    last_attempt_at = CASE WHEN ? THEN NOW() ELSE student_practice_progress.last_attempt_at END,
		    attempts_count = student_practice_progress.attempts_count + CASE WHEN ? THEN 1 ELSE 0 END,
		    approved_submission_id = COALESCE(EXCLUDED.approved_submission_id, student_practice_progress.approved_submission_id),
		    updated_at = NOW()
	`, uuid.New(),
		update.StudentID,
		update.PracticeID,
		update.CourseID,
		update.LessonID,
		update.Status,
		update.Status,
		update.IncrementRuns,
		update.ApprovedSubmissionID.String(),
		update.IncrementRuns,
		update.IncrementRuns,
	).Error
}

func syncUserLevelTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID) error {
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
