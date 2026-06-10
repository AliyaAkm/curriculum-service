package teacherstats

import (
	"context"
	teacherstatsdomain "curriculum-service/internal/domain/teacherstats"
	"math"
)

func (r *Repo) GetStatistics(ctx context.Context, filter teacherstatsdomain.Filter) (*teacherstatsdomain.Statistics, error) {
	summary, err := r.getSummary(ctx, filter)
	if err != nil {
		return nil, err
	}
	practice, err := r.getPracticeStats(ctx, filter)
	if err != nil {
		return nil, err
	}
	quiz, err := r.getQuizStats(ctx, filter)
	if err != nil {
		return nil, err
	}
	activity, err := r.getActivity(ctx, filter)
	if err != nil {
		return nil, err
	}
	funnel, err := r.getFunnel(ctx, filter)
	if err != nil {
		return nil, err
	}
	quizHeatmap, err := r.getQuizHeatmap(ctx, filter)
	if err != nil {
		return nil, err
	}
	courses, err := r.getCourseStats(ctx, filter)
	if err != nil {
		return nil, err
	}
	newStudents, err := r.getNewStudents(ctx, filter)
	if err != nil {
		return nil, err
	}
	queue, err := r.getReviewQueue(ctx, filter)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		activity = []teacherstatsdomain.ActivityDay{}
	}
	if funnel == nil {
		funnel = []teacherstatsdomain.FunnelStep{}
	}
	if quizHeatmap == nil {
		quizHeatmap = []teacherstatsdomain.QuizHeatmapItem{}
	}
	if courses == nil {
		courses = []teacherstatsdomain.CourseStats{}
	}
	if newStudents == nil {
		newStudents = []teacherstatsdomain.NewStudent{}
	}
	if queue == nil {
		queue = []teacherstatsdomain.ReviewItem{}
	}

	return &teacherstatsdomain.Statistics{
		Summary:     summary,
		Practice:    practice,
		Quiz:        quiz,
		Activity:    activity,
		Funnel:      funnel,
		QuizHeatmap: quizHeatmap,
		Courses:     courses,
		NewStudents: newStudents,
		ReviewQueue: queue,
	}, nil
}

func (r *Repo) getSummary(ctx context.Context, filter teacherstatsdomain.Filter) (teacherstatsdomain.Summary, error) {
	scope, args := teacherCourseScope(filter)
	query := `
		WITH ` + scope + `,
		subscriptions AS (
			SELECT
				cs.user_id,
				cs.course_id,
				MIN(cs.started_at) AS started_at,
				MAX(cs.last_activity_at) AS last_activity_at,
				MIN(cs.completed_at) AS completed_at
			FROM course_subscription cs
			INNER JOIN teacher_courses tc ON tc.id = cs.course_id
			GROUP BY cs.user_id, cs.course_id
		),
		lesson_counts AS (
			SELECT tc.id AS course_id, COUNT(DISTINCT cl.id)::int AS total_lessons
			FROM teacher_courses tc
			LEFT JOIN course_modules cm ON cm.course_id = tc.id
			LEFT JOIN course_lessons cl ON cl.module_id = cm.id
			GROUP BY tc.id
		),
		completed_lessons AS (
			SELECT
				s.user_id,
				s.course_id,
				COUNT(DISTINCT ucp.lesson_id)::int AS completed_lessons
			FROM subscriptions s
			LEFT JOIN course_modules cm ON cm.course_id = s.course_id
			LEFT JOIN course_lessons cl ON cl.module_id = cm.id
			LEFT JOIN user_course_points ucp ON ucp.lesson_id = cl.id AND ucp.user_id = s.user_id
			GROUP BY s.user_id, s.course_id
		),
		student_progress AS (
			SELECT
				s.user_id,
				s.course_id,
				s.started_at,
				s.completed_at,
				CASE
					WHEN lc.total_lessons > 0 THEN ROUND(COALESCE(clp.completed_lessons, 0)::numeric / lc.total_lessons * 100)
					ELSE 0
				END::int AS progress_percent
			FROM subscriptions s
			INNER JOIN lesson_counts lc ON lc.course_id = s.course_id
			LEFT JOIN completed_lessons clp ON clp.user_id = s.user_id AND clp.course_id = s.course_id
		),
		student_activity AS (
			-- started_at теперь = дата подписки (enrollment), а не первая активность.
			-- Поэтому в активность она НЕ входит, иначе только что записавшийся студент
			-- ошибочно считался бы "активным". Активность = last_activity_at + реальные действия.
			SELECT s.user_id, s.course_id, s.last_activity_at AS activity_at FROM subscriptions s
			UNION ALL
			SELECT ucp.user_id, cm.course_id, ucp.created_at AS activity_at
			FROM user_course_points ucp
			INNER JOIN course_lessons cl ON cl.id = ucp.lesson_id
			INNER JOIN course_modules cm ON cm.id = cl.module_id
			INNER JOIN teacher_courses tc ON tc.id = cm.course_id
			UNION ALL
			SELECT lqa.user_id, cm.course_id, lqa.updated_at AS activity_at
			FROM lesson_quiz_attempts lqa
			INNER JOIN lesson_quizzes lq ON lq.id = lqa.quiz_id
			INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
			INNER JOIN course_modules cm ON cm.id = cl.module_id
			INNER JOIN teacher_courses tc ON tc.id = cm.course_id
			UNION ALL
			SELECT cea.user_id, cea.course_id, cea.created_at AS activity_at
			FROM code_execution_attempts cea
			INNER JOIN teacher_courses tc ON tc.id = cea.course_id
			UNION ALL
			SELECT spp.student_id, spp.course_id, spp.last_attempt_at AS activity_at
			FROM student_practice_progress spp
			INNER JOIN teacher_courses tc ON tc.id = spp.course_id
			UNION ALL
			SELECT prs.student_id, prs.course_id, prs.updated_at AS activity_at
			FROM practice_review_submissions prs
			INNER JOIN teacher_courses tc ON tc.id = prs.course_id
		),
		student_last_activity AS (
			SELECT user_id, course_id, MAX(activity_at) AS last_activity_at
			FROM student_activity
			WHERE activity_at IS NOT NULL
			GROUP BY user_id, course_id
		)
		SELECT
			(SELECT COUNT(*) FROM teacher_courses)::int AS courses,
			COUNT(DISTINCT sp.user_id)::int AS students,
			COUNT(DISTINCT sp.user_id) FILTER (
				WHERE sp.started_at >= NOW() - (?::int * INTERVAL '1 day')
			)::int AS new_students,
			COUNT(DISTINCT sp.user_id) FILTER (
				WHERE sla.last_activity_at >= NOW() - (?::int * INTERVAL '1 day')
			)::int AS active_students,
			COUNT(DISTINCT sp.user_id) FILTER (WHERE sp.completed_at IS NOT NULL)::int AS completed_students,
			COALESCE(ROUND(AVG(sp.progress_percent)), 0)::int AS avg_progress_percent,
			COALESCE((
				SELECT SUM(uxe.xp)
				FROM user_xp_events uxe
				INNER JOIN teacher_courses tc ON tc.id = uxe.course_id
			), 0)::int AS total_xp_awarded,
			(SELECT MAX(last_activity_at) FROM student_last_activity) AS last_activity_at
		FROM student_progress sp
		LEFT JOIN student_last_activity sla ON sla.user_id = sp.user_id AND sla.course_id = sp.course_id
	`
	args = append(args, filter.PeriodDays, filter.PeriodDays)

	var row teacherstatsdomain.Summary
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&row).Error
	return row, err
}

func (r *Repo) getPracticeStats(ctx context.Context, filter teacherstatsdomain.Filter) (teacherstatsdomain.PracticeStats, error) {
	scope, args := teacherCourseScope(filter)
	query := `
		WITH ` + scope + `,
		auto_attempts AS (
			SELECT
				COUNT(*) FILTER (WHERE cea.run_type = 'run')::int AS auto_runs,
				COUNT(*) FILTER (WHERE cea.run_type = 'submit')::int AS auto_submissions,
				COUNT(*) FILTER (WHERE cea.run_type = 'submit' AND cea.passed = true)::int AS auto_successful_submits
			FROM code_execution_attempts cea
			INNER JOIN teacher_courses tc ON tc.id = cea.course_id
		),
		manual_submissions AS (
			SELECT
				COUNT(*)::int AS manual_submissions,
				COUNT(*) FILTER (WHERE prs.status = 'approved')::int AS manual_approved,
				COUNT(*) FILTER (WHERE prs.status IN ('submitted', 'in_review'))::int AS manual_pending_review,
				COUNT(*) FILTER (WHERE prs.status = 'changes_requested')::int AS manual_changes_requested
			FROM practice_review_submissions prs
			INNER JOIN teacher_courses tc ON tc.id = prs.course_id
		),
		progress AS (
			SELECT
				COUNT(DISTINCT (spp.student_id, spp.practice_id)) FILTER (WHERE spp.status = 'completed')::int AS completed_practices,
				COALESCE(ROUND(AVG(spp.attempts_count) FILTER (WHERE spp.status = 'completed')), 0)::int AS avg_attempts_per_completion
			FROM student_practice_progress spp
			INNER JOIN teacher_courses tc ON tc.id = spp.course_id
		)
		SELECT
			COALESCE(auto_attempts.auto_runs, 0)::int AS auto_runs,
			COALESCE(auto_attempts.auto_submissions, 0)::int AS auto_submissions,
			COALESCE(auto_attempts.auto_successful_submits, 0)::int AS auto_successful_submits,
			COALESCE(manual_submissions.manual_submissions, 0)::int AS manual_submissions,
			COALESCE(manual_submissions.manual_approved, 0)::int AS manual_approved,
			COALESCE(manual_submissions.manual_pending_review, 0)::int AS manual_pending_review,
			COALESCE(manual_submissions.manual_changes_requested, 0)::int AS manual_changes_requested,
			COALESCE(progress.completed_practices, 0)::int AS completed_practices,
			COALESCE(progress.avg_attempts_per_completion, 0)::int AS avg_attempts_per_completion
		FROM auto_attempts, manual_submissions, progress
	`

	var row teacherstatsdomain.PracticeStats
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&row).Error; err != nil {
		return teacherstatsdomain.PracticeStats{}, err
	}
	row.PassRatePercent = percent(row.AutoSuccessfulSubmits+row.ManualApproved, row.AutoSubmissions+row.ManualSubmissions)
	return row, nil
}

func (r *Repo) getQuizStats(ctx context.Context, filter teacherstatsdomain.Filter) (teacherstatsdomain.QuizStats, error) {
	scope, args := teacherCourseScope(filter)
	query := `
		WITH ` + scope + `
		SELECT
			COUNT(lqa.id)::int AS attempted_quizzes,
			COUNT(lqa.id) FILTER (WHERE lqa.is_correct = true)::int AS passed_quizzes
		FROM lesson_quiz_attempts lqa
		INNER JOIN lesson_quizzes lq ON lq.id = lqa.quiz_id
		INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		INNER JOIN teacher_courses tc ON tc.id = cm.course_id
	`

	var row teacherstatsdomain.QuizStats
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&row).Error; err != nil {
		return teacherstatsdomain.QuizStats{}, err
	}
	row.AccuracyPercent = percent(row.PassedQuizzes, row.AttemptedQuizzes)
	return row, nil
}

func (r *Repo) getActivity(ctx context.Context, filter teacherstatsdomain.Filter) ([]teacherstatsdomain.ActivityDay, error) {
	scope, args := teacherCourseScope(filter)
	query := `
		WITH ` + scope + `,
		days AS (
			SELECT generate_series(
				(CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day'))::date,
				CURRENT_DATE,
				INTERVAL '1 day'
			)::date AS day
		),
		subscriptions AS (
			SELECT DISTINCT cs.user_id, cs.course_id
			FROM course_subscription cs
			INNER JOIN teacher_courses tc ON tc.id = cs.course_id
		),
		student_activity AS (
			SELECT DISTINCT user_id, day
			FROM (
				-- started_at = дата подписки, а не ежедневная активность обучения,
				-- поэтому день подписки не учитывается как "активный день".
				SELECT s.user_id, cs.last_activity_at::date AS day
				FROM course_subscription cs
				INNER JOIN subscriptions s ON s.user_id = cs.user_id AND s.course_id = cs.course_id
				WHERE cs.last_activity_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
				UNION ALL
				SELECT ucp.user_id, ucp.created_at::date AS day
				FROM user_course_points ucp
				INNER JOIN course_lessons cl ON cl.id = ucp.lesson_id
				INNER JOIN course_modules cm ON cm.id = cl.module_id
				INNER JOIN subscriptions s ON s.user_id = ucp.user_id AND s.course_id = cm.course_id
				WHERE ucp.created_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
				UNION ALL
				SELECT lqa.user_id, lqa.updated_at::date AS day
				FROM lesson_quiz_attempts lqa
				INNER JOIN lesson_quizzes lq ON lq.id = lqa.quiz_id
				INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
				INNER JOIN course_modules cm ON cm.id = cl.module_id
				INNER JOIN subscriptions s ON s.user_id = lqa.user_id AND s.course_id = cm.course_id
				WHERE lqa.updated_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
				UNION ALL
				SELECT cea.user_id, cea.created_at::date AS day
				FROM code_execution_attempts cea
				INNER JOIN subscriptions s ON s.user_id = cea.user_id AND s.course_id = cea.course_id
				WHERE cea.created_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
				UNION ALL
				SELECT prs.student_id, prs.created_at::date AS day
				FROM practice_review_submissions prs
				INNER JOIN subscriptions s ON s.user_id = prs.student_id AND s.course_id = prs.course_id
				WHERE prs.created_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
			) source
			WHERE day IS NOT NULL
		),
		practice_submissions AS (
			SELECT day, SUM(total)::int AS total
			FROM (
				SELECT cea.created_at::date AS day, COUNT(*)::int AS total
				FROM code_execution_attempts cea
				INNER JOIN teacher_courses tc ON tc.id = cea.course_id
				WHERE cea.run_type = 'submit'
				  AND cea.created_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
				GROUP BY cea.created_at::date
				UNION ALL
				SELECT prs.created_at::date AS day, COUNT(*)::int AS total
				FROM practice_review_submissions prs
				INNER JOIN teacher_courses tc ON tc.id = prs.course_id
				WHERE prs.created_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
				GROUP BY prs.created_at::date
			) source
			GROUP BY day
		),
		practice_approved AS (
			SELECT day, SUM(total)::int AS total
			FROM (
				SELECT cea.created_at::date AS day, COUNT(*)::int AS total
				FROM code_execution_attempts cea
				INNER JOIN teacher_courses tc ON tc.id = cea.course_id
				WHERE cea.run_type = 'submit'
				  AND cea.passed = true
				  AND cea.created_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
				GROUP BY cea.created_at::date
				UNION ALL
				SELECT prs.reviewed_at::date AS day, COUNT(*)::int AS total
				FROM practice_review_submissions prs
				INNER JOIN teacher_courses tc ON tc.id = prs.course_id
				WHERE prs.status = 'approved'
				  AND prs.reviewed_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
				GROUP BY prs.reviewed_at::date
			) source
			GROUP BY day
		),
		quiz_attempts AS (
			SELECT lqa.updated_at::date AS day, COUNT(*)::int AS total
			FROM lesson_quiz_attempts lqa
			INNER JOIN lesson_quizzes lq ON lq.id = lqa.quiz_id
			INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
			INNER JOIN course_modules cm ON cm.id = cl.module_id
			INNER JOIN teacher_courses tc ON tc.id = cm.course_id
			WHERE lqa.updated_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
			GROUP BY lqa.updated_at::date
		),
		xp_events AS (
			SELECT uxe.activity_at::date AS day, COALESCE(SUM(uxe.xp), 0)::int AS total
			FROM user_xp_events uxe
			INNER JOIN teacher_courses tc ON tc.id = uxe.course_id
			WHERE uxe.activity_at >= CURRENT_DATE - ((?::int - 1) * INTERVAL '1 day')
			GROUP BY uxe.activity_at::date
		)
		SELECT
			TO_CHAR(days.day, 'YYYY-MM-DD') AS date,
			COALESCE(COUNT(DISTINCT student_activity.user_id), 0)::int AS active_students,
			COALESCE(practice_submissions.total, 0)::int AS practice_submissions,
			COALESCE(practice_approved.total, 0)::int AS practice_approved,
			COALESCE(quiz_attempts.total, 0)::int AS quiz_attempts,
			COALESCE(xp_events.total, 0)::int AS xp_awarded
		FROM days
		LEFT JOIN student_activity ON student_activity.day = days.day
		LEFT JOIN practice_submissions ON practice_submissions.day = days.day
		LEFT JOIN practice_approved ON practice_approved.day = days.day
		LEFT JOIN quiz_attempts ON quiz_attempts.day = days.day
		LEFT JOIN xp_events ON xp_events.day = days.day
		GROUP BY days.day, practice_submissions.total, practice_approved.total, quiz_attempts.total, xp_events.total
		ORDER BY days.day ASC
	`
	// 12 плейсхолдеров: days + last_activity + ucp + quiz + code + practice
	// + practice_submissions(cea,prs) + practice_approved(cea,prs) + quiz_attempts + xp_events.
	args = append(args,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
		filter.PeriodDays,
	)

	var rows []teacherstatsdomain.ActivityDay
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	return rows, err
}

func (r *Repo) getFunnel(ctx context.Context, filter teacherstatsdomain.Filter) ([]teacherstatsdomain.FunnelStep, error) {
	scope, args := teacherCourseScope(filter)
	query := `
		WITH ` + scope + `,
		subscriptions AS (
			SELECT
				cs.user_id,
				cs.course_id,
				MIN(cs.started_at) AS started_at,
				MAX(cs.last_activity_at) AS last_activity_at,
				MIN(cs.completed_at) AS completed_at
			FROM course_subscription cs
			INNER JOIN teacher_courses tc ON tc.id = cs.course_id
			GROUP BY cs.user_id, cs.course_id
		),
		lesson_completed AS (
			SELECT DISTINCT ucp.user_id, cm.course_id
			FROM user_course_points ucp
			INNER JOIN course_lessons cl ON cl.id = ucp.lesson_id
			INNER JOIN course_modules cm ON cm.id = cl.module_id
			INNER JOIN teacher_courses tc ON tc.id = cm.course_id
		),
		quiz_attempted AS (
			SELECT DISTINCT lqa.user_id, cm.course_id
			FROM lesson_quiz_attempts lqa
			INNER JOIN lesson_quizzes lq ON lq.id = lqa.quiz_id
			INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
			INNER JOIN course_modules cm ON cm.id = cl.module_id
			INNER JOIN teacher_courses tc ON tc.id = cm.course_id
		),
		activity AS (
			SELECT user_id, course_id, MAX(activity_at) AS last_activity_at
			FROM (
				-- started_at = дата подписки (enrollment), не активность: в воронке шаг
				-- "Started learning" должен отражать реальное обучение, а не сам факт подписки,
				-- иначе он всегда равнялся бы шагу "Enrolled".
				SELECT s.user_id, s.course_id, s.last_activity_at AS activity_at FROM subscriptions s
				UNION ALL
				SELECT lc.user_id, lc.course_id, ucp.created_at AS activity_at
				FROM lesson_completed lc
				INNER JOIN user_course_points ucp ON ucp.user_id = lc.user_id
				INNER JOIN course_lessons cl ON cl.id = ucp.lesson_id
				INNER JOIN course_modules cm ON cm.id = cl.module_id AND cm.course_id = lc.course_id
				UNION ALL
				SELECT qa.user_id, qa.course_id, lqa.updated_at AS activity_at
				FROM quiz_attempted qa
				INNER JOIN lesson_quiz_attempts lqa ON lqa.user_id = qa.user_id
				INNER JOIN lesson_quizzes lq ON lq.id = lqa.quiz_id
				INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
				INNER JOIN course_modules cm ON cm.id = cl.module_id AND cm.course_id = qa.course_id
				UNION ALL
				SELECT cea.user_id, cea.course_id, cea.created_at AS activity_at
				FROM code_execution_attempts cea
				INNER JOIN teacher_courses tc ON tc.id = cea.course_id
				UNION ALL
				SELECT prs.student_id AS user_id, prs.course_id, prs.updated_at AS activity_at
				FROM practice_review_submissions prs
				INNER JOIN teacher_courses tc ON tc.id = prs.course_id
			) source
			WHERE activity_at IS NOT NULL
			GROUP BY user_id, course_id
		)
		SELECT
			COUNT(*)::int AS enrolled,
			COUNT(*) FILTER (WHERE activity.user_id IS NOT NULL)::int AS started_learning,
			COUNT(*) FILTER (WHERE lesson_completed.user_id IS NOT NULL)::int AS completed_first_lesson,
			COUNT(*) FILTER (WHERE subscriptions.completed_at IS NOT NULL)::int AS completed_course
		FROM subscriptions
		LEFT JOIN activity ON activity.user_id = subscriptions.user_id AND activity.course_id = subscriptions.course_id
		LEFT JOIN lesson_completed ON lesson_completed.user_id = subscriptions.user_id AND lesson_completed.course_id = subscriptions.course_id
	`
	var row struct {
		Enrolled             int `gorm:"column:enrolled"`
		StartedLearning      int `gorm:"column:started_learning"`
		CompletedFirstLesson int `gorm:"column:completed_first_lesson"`
		CompletedCourse      int `gorm:"column:completed_course"`
	}
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&row).Error; err != nil {
		return nil, err
	}

	steps := []teacherstatsdomain.FunnelStep{
		{Key: "enrolled", Label: "Enrolled", Value: row.Enrolled},
		{Key: "started_learning", Label: "Started learning", Value: row.StartedLearning},
		{Key: "completed_first_lesson", Label: "Completed first lesson", Value: row.CompletedFirstLesson},
		{Key: "completed_course", Label: "Completed course", Value: row.CompletedCourse},
	}
	fillFunnelPercents(steps)
	return steps, nil
}

func (r *Repo) getQuizHeatmap(ctx context.Context, filter teacherstatsdomain.Filter) ([]teacherstatsdomain.QuizHeatmapItem, error) {
	scope, args := teacherCourseScope(filter)
	query := `
		WITH ` + scope + `
		SELECT
			lq.id AS quiz_id,
			tc.id AS course_id,
			tc.title AS course_title,
			cl.id AS lesson_id,
			COALESCE(lesson_title.name, '') AS lesson_title,
			COALESCE(lq.position, 0)::int AS position,
			COALESCE(NULLIF(quiz_text.question, ''), 'Quiz ' || COALESCE(lq.position, 0)::text) AS question,
			COUNT(lqa.id)::int AS attempts,
			COUNT(lqa.id) FILTER (WHERE lqa.is_correct = true)::int AS correct_attempts
		FROM lesson_quizzes lq
		INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
		INNER JOIN course_modules cm ON cm.id = cl.module_id
		INNER JOIN teacher_courses tc ON tc.id = cm.course_id
		LEFT JOIN lesson_quiz_attempts lqa ON lqa.quiz_id = lq.id
		LEFT JOIN LATERAL (
			SELECT lqt.question
			FROM lesson_quiz_texts lqt
			WHERE lqt.quiz_id = lq.id
			ORDER BY lqt.question
			LIMIT 1
		) quiz_text ON TRUE
		LEFT JOIN LATERAL (
			SELECT clt.name
			FROM course_lesson_titles clt
			WHERE clt.lesson_id = cl.id
			ORDER BY clt.name
			LIMIT 1
		) lesson_title ON TRUE
		GROUP BY
			lq.id,
			tc.id,
			tc.title,
			cl.id,
			lesson_title.name,
			lq.position,
			quiz_text.question,
			cm.position,
			cl.position
		ORDER BY
			tc.title ASC,
			COALESCE(cm.position, 0) ASC,
			COALESCE(cl.position, 0) ASC,
			COALESCE(lq.position, 0) ASC
		LIMIT 100
	`

	var rows []teacherstatsdomain.QuizHeatmapItem
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].AccuracyPercent = percent(rows[i].CorrectAttempts, rows[i].Attempts)
	}
	return rows, nil
}

func (r *Repo) getCourseStats(ctx context.Context, filter teacherstatsdomain.Filter) ([]teacherstatsdomain.CourseStats, error) {
	scope, args := teacherCourseScope(filter)
	query := `
		WITH ` + scope + `,
		subscriptions AS (
			SELECT
				cs.user_id,
				cs.course_id,
				MIN(cs.started_at) AS started_at,
				MAX(cs.last_activity_at) AS last_activity_at,
				MIN(cs.completed_at) AS completed_at
			FROM course_subscription cs
			INNER JOIN teacher_courses tc ON tc.id = cs.course_id
			GROUP BY cs.user_id, cs.course_id
		),
		lesson_counts AS (
			SELECT tc.id AS course_id, COUNT(DISTINCT cl.id)::int AS total_lessons
			FROM teacher_courses tc
			LEFT JOIN course_modules cm ON cm.course_id = tc.id
			LEFT JOIN course_lessons cl ON cl.module_id = cm.id
			GROUP BY tc.id
		),
		completed_lessons AS (
			SELECT
				s.user_id,
				s.course_id,
				COUNT(DISTINCT ucp.lesson_id)::int AS completed_lessons
			FROM subscriptions s
			LEFT JOIN course_modules cm ON cm.course_id = s.course_id
			LEFT JOIN course_lessons cl ON cl.module_id = cm.id
			LEFT JOIN user_course_points ucp ON ucp.lesson_id = cl.id AND ucp.user_id = s.user_id
			GROUP BY s.user_id, s.course_id
		),
		student_progress AS (
			SELECT
				s.user_id,
				s.course_id,
				s.completed_at,
				CASE
					WHEN lc.total_lessons > 0 THEN ROUND(COALESCE(clp.completed_lessons, 0)::numeric / lc.total_lessons * 100)
					ELSE 0
				END::int AS progress_percent
			FROM subscriptions s
			INNER JOIN lesson_counts lc ON lc.course_id = s.course_id
			LEFT JOIN completed_lessons clp ON clp.user_id = s.user_id AND clp.course_id = s.course_id
		),
		student_activity AS (
			-- started_at теперь = дата подписки (enrollment), а не первая активность.
			-- Поэтому в активность она НЕ входит, иначе только что записавшийся студент
			-- ошибочно считался бы "активным". Активность = last_activity_at + реальные действия.
			SELECT s.user_id, s.course_id, s.last_activity_at AS activity_at FROM subscriptions s
			UNION ALL
			SELECT ucp.user_id, cm.course_id, ucp.created_at AS activity_at
			FROM user_course_points ucp
			INNER JOIN course_lessons cl ON cl.id = ucp.lesson_id
			INNER JOIN course_modules cm ON cm.id = cl.module_id
			INNER JOIN teacher_courses tc ON tc.id = cm.course_id
			UNION ALL
			SELECT lqa.user_id, cm.course_id, lqa.updated_at AS activity_at
			FROM lesson_quiz_attempts lqa
			INNER JOIN lesson_quizzes lq ON lq.id = lqa.quiz_id
			INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
			INNER JOIN course_modules cm ON cm.id = cl.module_id
			INNER JOIN teacher_courses tc ON tc.id = cm.course_id
			UNION ALL
			SELECT cea.user_id, cea.course_id, cea.created_at AS activity_at
			FROM code_execution_attempts cea
			INNER JOIN teacher_courses tc ON tc.id = cea.course_id
			UNION ALL
			SELECT spp.student_id, spp.course_id, spp.last_attempt_at AS activity_at
			FROM student_practice_progress spp
			INNER JOIN teacher_courses tc ON tc.id = spp.course_id
			UNION ALL
			SELECT prs.student_id, prs.course_id, prs.updated_at AS activity_at
			FROM practice_review_submissions prs
			INNER JOIN teacher_courses tc ON tc.id = prs.course_id
		),
		student_last_activity AS (
			SELECT user_id, course_id, MAX(activity_at) AS last_activity_at
			FROM student_activity
			WHERE activity_at IS NOT NULL
			GROUP BY user_id, course_id
		),
		practice_queue AS (
			SELECT course_id, COUNT(*) FILTER (WHERE status IN ('submitted', 'in_review'))::int AS pending_review
			FROM practice_review_submissions
			GROUP BY course_id
		),
		practice_completed AS (
			SELECT course_id, COUNT(DISTINCT (student_id, practice_id)) FILTER (WHERE status = 'completed')::int AS completed
			FROM student_practice_progress
			GROUP BY course_id
		),
		quiz_stats AS (
			SELECT
				cm.course_id,
				COUNT(lqa.id)::int AS attempted,
				COUNT(lqa.id) FILTER (WHERE lqa.is_correct = true)::int AS passed
			FROM lesson_quiz_attempts lqa
			INNER JOIN lesson_quizzes lq ON lq.id = lqa.quiz_id
			INNER JOIN course_lessons cl ON cl.id = lq.lesson_id
			INNER JOIN course_modules cm ON cm.id = cl.module_id
			GROUP BY cm.course_id
		),
		xp AS (
			SELECT course_id, COALESCE(SUM(xp), 0)::int AS awarded
			FROM user_xp_events
			GROUP BY course_id
		)
		SELECT
			tc.id AS course_id,
			tc.title,
			COUNT(DISTINCT sp.user_id)::int AS students,
			COUNT(DISTINCT sp.user_id) FILTER (
				WHERE sla.last_activity_at >= NOW() - (?::int * INTERVAL '1 day')
			)::int AS active_students,
			COUNT(DISTINCT sp.user_id) FILTER (WHERE sp.completed_at IS NOT NULL)::int AS completed_students,
			COALESCE(lc.total_lessons, 0)::int AS total_lessons,
			COALESCE(ROUND(AVG(sp.progress_percent)), 0)::int AS avg_progress_percent,
			COALESCE(pq.pending_review, 0)::int AS practice_pending_review,
			COALESCE(pc.completed, 0)::int AS practice_completed,
			COALESCE(
				ROUND(CASE WHEN qs.attempted > 0 THEN qs.passed::numeric / qs.attempted * 100 ELSE 0 END),
				0
			)::int AS quiz_accuracy_percent,
			COALESCE(xp.awarded, 0)::int AS xp_awarded,
			MAX(sla.last_activity_at) AS last_activity_at
		FROM teacher_courses tc
		LEFT JOIN student_progress sp ON sp.course_id = tc.id
		LEFT JOIN student_last_activity sla ON sla.user_id = sp.user_id AND sla.course_id = sp.course_id
		LEFT JOIN lesson_counts lc ON lc.course_id = tc.id
		LEFT JOIN practice_queue pq ON pq.course_id = tc.id
		LEFT JOIN practice_completed pc ON pc.course_id = tc.id
		LEFT JOIN quiz_stats qs ON qs.course_id = tc.id
		LEFT JOIN xp ON xp.course_id = tc.id
		GROUP BY tc.id, tc.title, lc.total_lessons, pq.pending_review, pc.completed, qs.attempted, qs.passed, xp.awarded
		ORDER BY MAX(sla.last_activity_at) DESC NULLS LAST, tc.title ASC
	`
	args = append(args, filter.PeriodDays)

	var rows []teacherstatsdomain.CourseStats
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	return rows, err
}
func (r *Repo) getNewStudents(ctx context.Context, filter teacherstatsdomain.Filter) ([]teacherstatsdomain.NewStudent, error) {
	scope, args := teacherCourseScope(filter)
	query := `
		WITH ` + scope + `,
		recent_subscriptions AS (
			-- одна строка на студента: его самый свежий старт курса преподавателя
			SELECT DISTINCT ON (cs.user_id)
				cs.user_id,
				cs.course_id,
				cs.started_at AS joined_at
			FROM course_subscription cs
			INNER JOIN teacher_courses tc ON tc.id = cs.course_id
			WHERE cs.started_at IS NOT NULL
			  AND cs.started_at >= NOW() - (?::int * INTERVAL '1 day')
			ORDER BY cs.user_id, cs.started_at DESC
		)
		SELECT
			rs.user_id AS student_id,
			COALESCE(u.email, '') AS student_email,
			COALESCE(NULLIF(u.login, ''), split_part(COALESCE(u.email, ''), '@', 1), '') AS student_name,
			COALESCE(u.photo_url, '') AS photo_url,
			rs.course_id,
			tc.title AS course_title,
			rs.joined_at AS subscribed_at
		FROM recent_subscriptions rs
		INNER JOIN teacher_courses tc ON tc.id = rs.course_id
		LEFT JOIN users u ON u.id = rs.user_id
		ORDER BY rs.joined_at DESC, tc.title ASC
		LIMIT 20
	`
	args = append(args, filter.PeriodDays)

	var rows []teacherstatsdomain.NewStudent
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	return rows, err
}

func (r *Repo) getReviewQueue(ctx context.Context, filter teacherstatsdomain.Filter) ([]teacherstatsdomain.ReviewItem, error) {
	scope, args := teacherCourseScope(filter)
	query := `
		WITH ` + scope + `
		SELECT
			prs.id AS submission_id,
			prs.course_id,
			tc.title AS course_title,
			prs.practice_id,
			pt.title AS practice_title,
			prs.student_id,
			COALESCE(u.email, '') AS student_email,
			prs.status,
			prs.attempt_number,
			prs.created_at
		FROM practice_review_submissions prs
		INNER JOIN teacher_courses tc ON tc.id = prs.course_id
		INNER JOIN practice_tasks pt ON pt.id = prs.practice_id
		LEFT JOIN users u ON u.id = prs.student_id
		WHERE prs.status IN ('submitted', 'in_review')
		ORDER BY
			CASE prs.status WHEN 'submitted' THEN 0 ELSE 1 END,
			prs.created_at ASC
		LIMIT 20
	`

	var rows []teacherstatsdomain.ReviewItem
	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	return rows, err
}
func teacherCourseScope(filter teacherstatsdomain.Filter) (string, []any) {
	args := make([]any, 0, 2)
	query := `
		teacher_courses AS (
			SELECT c.id, c.title
			FROM courses c
			WHERE 1 = 1
	`

	// Админ видит все курсы; обычный преподаватель — только свои.
	if !filter.IsAdmin {
		// ВАРИАНТ A: если у courses есть колонка автора/преподавателя.
		query += `
			  AND c.author_id = ?
		`
		args = append(args, filter.TeacherID)

		// ВАРИАНТ B: если связь через таблицу course_teachers — используйте это вместо A:
		// query += `
		//   AND EXISTS (
		//       SELECT 1 FROM course_teachers ct
		//       WHERE ct.course_id = c.id AND ct.teacher_id = ?
		//   )
		// `
		// args = append(args, filter.TeacherID)
	}

	if filter.CourseID != nil {
		query += `
			  AND c.id = ?
		`
		args = append(args, *filter.CourseID)
	}

	query += `
		)
	`
	return query, args
}
func percent(done, total int) int {
	if total <= 0 {
		return 0
	}
	return int(math.Round(float64(done) / float64(total) * 100))
}

func fillFunnelPercents(steps []teacherstatsdomain.FunnelStep) {
	if len(steps) == 0 {
		return
	}
	for i := 1; i < len(steps); i++ {
		if steps[i].Value > steps[i-1].Value {
			steps[i].Value = steps[i-1].Value
		}
	}

	total := steps[0].Value
	previous := total

	for i := range steps {
		steps[i].ConversionPercent = percent(steps[i].Value, total)
		if i == 0 {
			steps[i].DropOffPercent = 0
			continue
		}
		drop := previous - steps[i].Value
		if drop < 0 {
			drop = 0
		}
		steps[i].DropOffPercent = percent(drop, previous)
		previous = steps[i].Value
	}
}
