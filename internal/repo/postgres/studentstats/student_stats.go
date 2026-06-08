package studentstats

import (
	"context"
	"curriculum-service/internal/domain/studentstats"
	"math"

	"github.com/google/uuid"
)

func (r *Repo) GetStatistics(ctx context.Context, userID uuid.UUID) (*studentstats.Statistics, error) {
	summary, err := r.getSummary(ctx, userID)
	if err != nil {
		return nil, err
	}
	quiz, err := r.getQuizStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	practice, err := r.getPracticeStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	activity, err := r.getActivity(ctx, userID)
	if err != nil {
		return nil, err
	}
	topics, err := r.getTopicProgress(ctx, userID)
	if err != nil {
		return nil, err
	}
	courses, err := r.getCourseProgress(ctx, userID)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		activity = []studentstats.ActivityDay{}
	}
	if topics == nil {
		topics = []studentstats.TopicProgress{}
	}
	if courses == nil {
		courses = []studentstats.CourseDetail{}
	}

	return &studentstats.Statistics{
		Summary:        summary,
		Quiz:           quiz,
		Practice:       practice,
		Activity:       activity,
		Topics:         topics,
		CourseProgress: courses,
	}, nil
}

func (r *Repo) getSummary(ctx context.Context, userID uuid.UUID) (studentstats.Summary, error) {
	var row studentstats.Summary
	row.UserID = userID

	err := r.db.WithContext(ctx).Raw(`
		WITH last_activity AS (
			SELECT MAX(activity_at) AS value
			FROM (
				SELECT MAX(activity_at) AS activity_at FROM user_xp_events WHERE user_id = ?
				UNION ALL
				SELECT MAX(updated_at) AS activity_at FROM lesson_quiz_attempts WHERE user_id = ?
				UNION ALL
				SELECT MAX(created_at) AS activity_at FROM code_execution_attempts WHERE user_id = ?
				UNION ALL
				SELECT MAX(last_attempt_at) AS activity_at FROM student_practice_progress WHERE student_id = ?
				UNION ALL
				SELECT MAX(last_activity_at) AS activity_at FROM course_subscription WHERE user_id = ?
			) source
		)
		SELECT
			COALESCE((SELECT SUM(xp) FROM user_xp_events WHERE user_id = ?), 0)::int AS total_xp,
			GREATEST(
				COALESCE((SELECT level FROM users WHERE id = ?), 1),
				1 + (
					COALESCE((SELECT SUM(xp) FROM user_xp_events WHERE user_id = ?), 0)::bigint / 180
				)
			)::int AS level,
			COALESCE((SELECT streak FROM daily_streak WHERE user_id = ?), 0)::int AS current_streak,
			COALESCE((SELECT max_streak FROM users WHERE id = ?), 0)::int AS max_streak,
			COALESCE((SELECT COUNT(*) FROM course_subscription WHERE user_id = ? AND started_at IS NOT NULL), 0)::int AS started_courses,
			COALESCE((SELECT COUNT(*) FROM course_subscription WHERE user_id = ? AND started_at IS NOT NULL AND completed_at IS NULL), 0)::int AS active_courses,
			COALESCE((SELECT COUNT(*) FROM course_subscription WHERE user_id = ? AND completed_at IS NOT NULL), 0)::int AS completed_courses,
			COALESCE((
				SELECT COUNT(DISTINCT cl.id)
				FROM course_subscription cs
				INNER JOIN course_modules cm ON cm.course_id = cs.course_id
				INNER JOIN course_lessons cl ON cl.module_id = cm.id
				WHERE cs.user_id = ?
			), 0)::int AS total_lessons,
			COALESCE((SELECT COUNT(DISTINCT lesson_id) FROM user_course_points WHERE user_id = ?), 0)::int AS completed_lessons,
			COALESCE((SELECT COUNT(*) FROM course_certificates WHERE user_id = ?), 0)::int AS certificates,
			COALESCE((SELECT COUNT(*) FROM user_achievements WHERE user_id = ? AND unlocked = true), 0)::int AS achievements,
			COALESCE((SELECT COUNT(*) FROM achievements WHERE is_active = true), 0)::int AS total_achievements,
			last_activity.value AS last_activity_at
		FROM last_activity
	`, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID).Scan(&row).Error
	if err != nil {
		return studentstats.Summary{}, err
	}
	row.ProgressPercent = percent(row.CompletedLessons, row.TotalLessons)
	return row, nil
}

func (r *Repo) getQuizStats(ctx context.Context, userID uuid.UUID) (studentstats.QuizStats, error) {
	var row studentstats.QuizStats
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(*)::int AS attempted_quizzes,
			COUNT(*) FILTER (WHERE is_correct = true)::int AS passed_quizzes
		FROM lesson_quiz_attempts
		WHERE user_id = ?
	`, userID).Scan(&row).Error
	if err != nil {
		return studentstats.QuizStats{}, err
	}
	row.AccuracyPercent = percent(row.PassedQuizzes, row.AttemptedQuizzes)
	return row, nil
}

func (r *Repo) getPracticeStats(ctx context.Context, userID uuid.UUID) (studentstats.PracticeStats, error) {
	var row studentstats.PracticeStats
	err := r.db.WithContext(ctx).Raw(`
		WITH auto_attempts AS (
			SELECT
				COUNT(*) FILTER (WHERE run_type = 'run')::int AS runs,
				COUNT(*) FILTER (WHERE run_type = 'submit')::int AS auto_submissions,
				COUNT(*) FILTER (WHERE run_type = 'submit' AND passed = true)::int AS successful_submits,
				COUNT(DISTINCT practice_id)::int AS attempted_practices,
				MAX(created_at) AS last_attempt_at
			FROM code_execution_attempts
			WHERE user_id = ?
		),
		manual_submissions AS (
			SELECT
				COUNT(*)::int AS submissions,
				COUNT(*) FILTER (WHERE status = 'approved')::int AS approved_submissions,
				COUNT(*) FILTER (WHERE status IN ('submitted', 'in_review'))::int AS pending_reviews,
				COUNT(*) FILTER (WHERE status = 'changes_requested')::int AS changes_requested,
				MAX(created_at) AS last_attempt_at
			FROM practice_review_submissions
			WHERE student_id = ?
		),
		progress AS (
			SELECT
				COUNT(DISTINCT practice_id)::int AS attempted_practices,
				COUNT(DISTINCT practice_id) FILTER (WHERE status = 'completed')::int AS completed_practices,
				MAX(last_attempt_at) AS last_attempt_at
			FROM student_practice_progress
			WHERE student_id = ?
		)
		SELECT
			COALESCE(auto_attempts.runs, 0)::int AS runs,
			(COALESCE(auto_attempts.auto_submissions, 0) + COALESCE(manual_submissions.submissions, 0))::int AS submissions,
			(COALESCE(auto_attempts.successful_submits, 0) + COALESCE(manual_submissions.approved_submissions, 0))::int AS successful_submits,
			GREATEST(COALESCE(auto_attempts.attempted_practices, 0), COALESCE(progress.attempted_practices, 0))::int AS attempted_practices,
			COALESCE(progress.completed_practices, 0)::int AS completed_practices,
			COALESCE((SELECT SUM(xp) FROM practice_xp_awards WHERE user_id = ?), 0)::int AS xp_earned,
			COALESCE(auto_attempts.runs, 0)::int AS auto_runs,
			COALESCE(auto_attempts.auto_submissions, 0)::int AS auto_submissions,
			COALESCE(auto_attempts.successful_submits, 0)::int AS auto_successful_submits,
			COALESCE(manual_submissions.submissions, 0)::int AS manual_submissions,
			COALESCE(manual_submissions.approved_submissions, 0)::int AS manual_approved,
			COALESCE(manual_submissions.pending_reviews, 0)::int AS manual_pending_review,
			COALESCE(manual_submissions.changes_requested, 0)::int AS manual_changes_requested,
			(
				SELECT MAX(value)
				FROM (VALUES
					(auto_attempts.last_attempt_at),
					(manual_submissions.last_attempt_at),
					(progress.last_attempt_at)
				) AS activity(value)
			) AS last_attempt_at
		FROM auto_attempts, manual_submissions, progress
	`, userID, userID, userID, userID).Scan(&row).Error
	if err != nil {
		return studentstats.PracticeStats{}, err
	}
	row.PassRatePercent = percent(row.SuccessfulSubmits, row.Submissions)
	return row, nil
}

func (r *Repo) getActivity(ctx context.Context, userID uuid.UUID) ([]studentstats.ActivityDay, error) {
	var rows []studentstats.ActivityDay
	err := r.db.WithContext(ctx).Raw(`
		WITH days AS (
			SELECT generate_series(
				(CURRENT_DATE - INTERVAL '29 days')::date,
				CURRENT_DATE,
				INTERVAL '1 day'
			)::date AS day
		),
		lessons AS (
			SELECT created_at::date AS day, COUNT(*)::int AS lessons_completed
			FROM user_course_points
			WHERE user_id = ?
			  AND created_at >= CURRENT_DATE - INTERVAL '29 days'
			GROUP BY created_at::date
		),
		quizzes AS (
			SELECT updated_at::date AS day, COUNT(*) FILTER (WHERE is_correct = true)::int AS quizzes_passed
			FROM lesson_quiz_attempts
			WHERE user_id = ?
			  AND updated_at >= CURRENT_DATE - INTERVAL '29 days'
			GROUP BY updated_at::date
		),
		practice AS (
			SELECT
				day,
				SUM(practice_attempts)::int AS practice_attempts,
				SUM(practice_completed)::int AS practice_completed
			FROM (
				SELECT
					created_at::date AS day,
					COUNT(*)::int AS practice_attempts,
					0::int AS practice_completed
				FROM code_execution_attempts
				WHERE user_id = ?
				  AND created_at >= CURRENT_DATE - INTERVAL '29 days'
				GROUP BY created_at::date

				UNION ALL

				SELECT
					created_at::date AS day,
					COUNT(*)::int AS practice_attempts,
					0::int AS practice_completed
				FROM practice_review_submissions
				WHERE student_id = ?
				  AND created_at >= CURRENT_DATE - INTERVAL '29 days'
				GROUP BY created_at::date

				UNION ALL

				SELECT
					completed_at::date AS day,
					0::int AS practice_attempts,
					COUNT(*)::int AS practice_completed
				FROM student_practice_progress
				WHERE student_id = ?
				  AND status = 'completed'
				  AND completed_at >= CURRENT_DATE - INTERVAL '29 days'
				GROUP BY completed_at::date
			) source
			GROUP BY day
		),
		xp_events AS (
			SELECT activity_at::date AS day, COALESCE(SUM(xp), 0)::int AS xp
			FROM user_xp_events
			WHERE user_id = ?
			  AND activity_at >= CURRENT_DATE - INTERVAL '29 days'
			GROUP BY activity_at::date
		)
		SELECT
			TO_CHAR(days.day, 'YYYY-MM-DD') AS date,
			COALESCE(lessons.lessons_completed, 0)::int AS lessons_completed,
			COALESCE(quizzes.quizzes_passed, 0)::int AS quizzes_passed,
			COALESCE(practice.practice_attempts, 0)::int AS practice_attempts,
			COALESCE(practice.practice_completed, 0)::int AS practice_completed,
			0::int AS ai_requests,
			COALESCE(xp_events.xp, 0)::int AS xp
		FROM days
		LEFT JOIN lessons ON lessons.day = days.day
		LEFT JOIN quizzes ON quizzes.day = days.day
		LEFT JOIN practice ON practice.day = days.day
		LEFT JOIN xp_events ON xp_events.day = days.day
		ORDER BY days.day ASC
	`, userID, userID, userID, userID, userID, userID).Scan(&rows).Error
	return rows, err
}

func (r *Repo) getTopicProgress(ctx context.Context, userID uuid.UUID) ([]studentstats.TopicProgress, error) {
	var rows []studentstats.TopicProgress
	err := r.db.WithContext(ctx).Raw(`
		WITH lesson_progress AS (
			SELECT
				ct.id AS topic_id,
				ct.code,
				ct.name,
				COUNT(DISTINCT cl.id)::int AS total_lessons,
				COUNT(DISTINCT ucp.lesson_id)::int AS completed_lessons,
				COALESCE(SUM(ucp.xp), 0)::int AS lesson_xp
			FROM course_subscription cs
			INNER JOIN courses c ON c.id = cs.course_id
			INNER JOIN course_topics ct ON ct.id = c.topic_id
			INNER JOIN course_modules cm ON cm.course_id = c.id
			INNER JOIN course_lessons cl ON cl.module_id = cm.id
			LEFT JOIN user_course_points ucp ON ucp.lesson_id = cl.id AND ucp.user_id = cs.user_id
			WHERE cs.user_id = ?
			GROUP BY ct.id, ct.code, ct.name
		),
		xp_progress AS (
			SELECT
				ct.id AS topic_id,
				COALESCE(SUM(uxe.xp), 0)::int AS xp
			FROM user_xp_events uxe
			INNER JOIN courses c ON c.id = uxe.course_id
			INNER JOIN course_topics ct ON ct.id = c.topic_id
			WHERE uxe.user_id = ?
			GROUP BY ct.id
		)
		SELECT
			lesson_progress.topic_id,
			lesson_progress.code,
			lesson_progress.name,
			lesson_progress.total_lessons,
			lesson_progress.completed_lessons,
			COALESCE(xp_progress.xp, 0)::int AS xp
		FROM lesson_progress
		LEFT JOIN xp_progress ON xp_progress.topic_id = lesson_progress.topic_id
		ORDER BY completed_lessons DESC, total_lessons DESC, lesson_progress.name ASC
	`, userID, userID).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].ProgressPercent = percent(rows[i].CompletedLessons, rows[i].TotalLessons)
	}
	return rows, nil
}

func (r *Repo) getCourseProgress(ctx context.Context, userID uuid.UUID) ([]studentstats.CourseDetail, error) {
	var courseRows []studentstats.CourseDetail
	if err := r.db.WithContext(ctx).Raw(`
		WITH subscriptions AS (
			SELECT
				user_id,
				course_id,
				MIN(started_at) AS started_at,
				MAX(last_activity_at) AS last_activity_at,
				MIN(completed_at) AS completed_at,
				(
					ARRAY_AGG(current_lesson_id ORDER BY last_activity_at DESC NULLS LAST, started_at DESC NULLS LAST)
					FILTER (WHERE current_lesson_id IS NOT NULL)
				)[1] AS current_lesson_id
			FROM course_subscription
			WHERE user_id = ?
			GROUP BY user_id, course_id
		)
		SELECT
			c.id AS course_id,
			c.title,
			cs.started_at,
			(
				SELECT MAX(value)
				FROM (VALUES
					(cs.last_activity_at),
					(MAX(cea.created_at)),
					(MAX(spp.last_attempt_at)),
					(MAX(prs.updated_at))
				) AS activity(value)
			) AS last_activity_at,
			cs.completed_at,
			cs.current_lesson_id,
			COUNT(DISTINCT cl.id)::int AS total_lessons,
			COUNT(DISTINCT ucp.lesson_id)::int AS completed_lessons,
			COUNT(DISTINCT cea.id) FILTER (WHERE cea.run_type = 'run')::int AS practice_runs,
			(
				COUNT(DISTINCT cea.id) FILTER (WHERE cea.run_type = 'submit') +
				COUNT(DISTINCT prs.id)
			)::int AS practice_submissions,
			COUNT(DISTINCT spp.practice_id) FILTER (WHERE spp.status = 'completed')::int AS practice_completed,
			COUNT(DISTINCT prs.id) FILTER (WHERE prs.status IN ('submitted', 'in_review'))::int AS practice_pending_review,
			COUNT(DISTINCT prs.id) FILTER (WHERE prs.status = 'changes_requested')::int AS practice_changes_requested,
			COALESCE((
				SELECT SUM(pxa.xp)
				FROM practice_xp_awards pxa
				WHERE pxa.user_id = cs.user_id
				  AND pxa.course_id = c.id
			), 0)::int AS practice_xp_earned
		FROM subscriptions cs
		INNER JOIN courses c ON c.id = cs.course_id
		LEFT JOIN course_modules cm ON cm.course_id = c.id
		LEFT JOIN course_lessons cl ON cl.module_id = cm.id
		LEFT JOIN user_course_points ucp ON ucp.lesson_id = cl.id AND ucp.user_id = cs.user_id
		LEFT JOIN code_execution_attempts cea ON cea.course_id = c.id AND cea.user_id = cs.user_id
		LEFT JOIN student_practice_progress spp ON spp.course_id = c.id AND spp.student_id = cs.user_id
		LEFT JOIN practice_review_submissions prs ON prs.course_id = c.id AND prs.student_id = cs.user_id
		GROUP BY c.id, c.title, cs.user_id, cs.started_at, cs.last_activity_at, cs.completed_at, cs.current_lesson_id
		ORDER BY COALESCE((
			SELECT MAX(value)
			FROM (VALUES
				(cs.last_activity_at),
				(MAX(cea.created_at)),
				(MAX(spp.last_attempt_at)),
				(MAX(prs.updated_at)),
				(cs.started_at),
				(c.created_at)
			) AS activity(value)
		), c.created_at) DESC
	`, userID).Scan(&courseRows).Error; err != nil {
		return nil, err
	}
	for i := range courseRows {
		courseRows[i].ProgressPercent = percent(courseRows[i].CompletedLessons, courseRows[i].TotalLessons)
		courseRows[i].Modules = []studentstats.ModuleDetail{}
	}

	var moduleRows []struct {
		CourseID         uuid.UUID `gorm:"column:course_id"`
		ModuleID         uuid.UUID `gorm:"column:module_id"`
		Title            string    `gorm:"column:title"`
		Position         int       `gorm:"column:position"`
		TotalLessons     int       `gorm:"column:total_lessons"`
		CompletedLessons int       `gorm:"column:completed_lessons"`
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			cm.course_id,
			cm.id AS module_id,
			cm.title,
			COALESCE(cm.position, 0)::int AS position,
			COUNT(DISTINCT cl.id)::int AS total_lessons,
			COUNT(DISTINCT ucp.lesson_id)::int AS completed_lessons
		FROM course_subscription cs
		INNER JOIN course_modules cm ON cm.course_id = cs.course_id
		LEFT JOIN course_lessons cl ON cl.module_id = cm.id
		LEFT JOIN user_course_points ucp ON ucp.lesson_id = cl.id AND ucp.user_id = cs.user_id
		WHERE cs.user_id = ?
		GROUP BY cm.course_id, cm.id, cm.title, cm.position, cm.created_at
		ORDER BY cm.course_id, COALESCE(cm.position, 0), cm.created_at
	`, userID).Scan(&moduleRows).Error; err != nil {
		return nil, err
	}

	byCourse := make(map[uuid.UUID]int, len(courseRows))
	for i := range courseRows {
		byCourse[courseRows[i].CourseID] = i
	}

	openByCourse := make(map[uuid.UUID]bool, len(courseRows))
	for _, course := range courseRows {
		openByCourse[course.CourseID] = true
	}
	for _, row := range moduleRows {
		idx, ok := byCourse[row.CourseID]
		if !ok {
			continue
		}
		isOpen := openByCourse[row.CourseID]
		courseRows[idx].Modules = append(courseRows[idx].Modules, studentstats.ModuleDetail{
			ModuleID:         row.ModuleID,
			Title:            row.Title,
			Position:         row.Position,
			TotalLessons:     row.TotalLessons,
			CompletedLessons: row.CompletedLessons,
			ProgressPercent:  percent(row.CompletedLessons, row.TotalLessons),
			IsOpen:           isOpen,
		})
		if row.TotalLessons > 0 && row.CompletedLessons < row.TotalLessons {
			openByCourse[row.CourseID] = false
		}
	}

	return courseRows, nil
}

func percent(done, total int) int {
	if total <= 0 {
		return 0
	}
	return int(math.Round(float64(done) / float64(total) * 100))
}
