package achievement

import (
	"context"
	achievementdomain "curriculum-service/internal/domain/achievement"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type achievementRow struct {
	ID            uuid.UUID  `gorm:"column:id"`
	Code          string     `gorm:"column:code"`
	IconKey       string     `gorm:"column:icon_key"`
	MetricKey     string     `gorm:"column:metric_key"`
	GoalValue     int        `gorm:"column:goal_value"`
	SortOrder     int        `gorm:"column:sort_order"`
	TitleEN       string     `gorm:"column:title_en"`
	TitleRU       string     `gorm:"column:title_ru"`
	TitleKK       string     `gorm:"column:title_kk"`
	DescriptionEN string     `gorm:"column:description_en"`
	DescriptionRU string     `gorm:"column:description_ru"`
	DescriptionKK string     `gorm:"column:description_kk"`
	ProgressValue int        `gorm:"column:progress_value"`
	Unlocked      bool       `gorm:"column:unlocked"`
	UnlockedAt    *time.Time `gorm:"column:unlocked_at"`
}

type userAchievementRow struct {
	AchievementID uuid.UUID  `gorm:"column:achievement_id"`
	UnlockedAt    *time.Time `gorm:"column:unlocked_at"`
}

type metricTotals struct {
	CompletedLessons int `gorm:"column:completed_lessons"`
	PassedQuizzes    int `gorm:"column:passed_quizzes"`
	MaxStreak        int `gorm:"column:max_streak"`
	XPTotal          int `gorm:"column:xp_total"`
	StartedCourses   int `gorm:"column:started_courses"`
	CompletedCourses int `gorm:"column:completed_courses"`
	ActiveCourses    int `gorm:"column:active_courses"`
	CompletedModules int `gorm:"column:completed_modules"`
}

func (r *Repo) ListAchievements(ctx context.Context, userID uuid.UUID) ([]achievementdomain.Achievement, error) {
	if err := r.syncUserAchievements(ctx, userID); err != nil {
		return nil, err
	}

	var rows []achievementRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT
			a.id,
			a.code,
			a.icon_key,
			a.metric_key,
			a.goal_value,
			a.sort_order,
			COALESCE(MAX(at.title) FILTER (WHERE cl.code = 'en'), '') AS title_en,
			COALESCE(MAX(at.title) FILTER (WHERE cl.code = 'ru'), '') AS title_ru,
			COALESCE(MAX(at.title) FILTER (WHERE cl.code = 'kk'), '') AS title_kk,
			COALESCE(MAX(at.description) FILTER (WHERE cl.code = 'en'), '') AS description_en,
			COALESCE(MAX(at.description) FILTER (WHERE cl.code = 'ru'), '') AS description_ru,
			COALESCE(MAX(at.description) FILTER (WHERE cl.code = 'kk'), '') AS description_kk,
			COALESCE(ua.progress_value, 0)::int AS progress_value,
			COALESCE(ua.unlocked, false) AS unlocked,
			ua.unlocked_at
		FROM achievements a
		LEFT JOIN achievement_texts at ON at.achievement_id = a.id
		LEFT JOIN course_locales cl ON cl.id = at.locale_id
		LEFT JOIN user_achievements ua ON ua.achievement_id = a.id AND ua.user_id = ?
		WHERE a.is_active = true
		GROUP BY
			a.id,
			a.code,
			a.icon_key,
			a.metric_key,
			a.goal_value,
			a.sort_order,
			ua.progress_value,
			ua.unlocked,
			ua.unlocked_at
		ORDER BY a.sort_order ASC, a.created_at ASC
	`, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make([]achievementdomain.Achievement, len(rows))
	for i := range rows {
		items[i] = achievementdomain.Achievement{
			ID:   rows[i].ID,
			Code: rows[i].Code,
			Title: achievementdomain.LocalizedText{
				EN: rows[i].TitleEN,
				RU: rows[i].TitleRU,
				KK: rows[i].TitleKK,
			},
			Description: achievementdomain.LocalizedText{
				EN: rows[i].DescriptionEN,
				RU: rows[i].DescriptionRU,
				KK: rows[i].DescriptionKK,
			},
			IconKey:    rows[i].IconKey,
			MetricKey:  rows[i].MetricKey,
			Goal:       rows[i].GoalValue,
			Progress:   min(rows[i].ProgressValue, rows[i].GoalValue),
			Unlocked:   rows[i].Unlocked,
			UnlockedAt: rows[i].UnlockedAt,
			SortOrder:  rows[i].SortOrder,
		}
	}

	return items, nil
}

func (r *Repo) syncUserAchievements(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var achievements []struct {
			ID        uuid.UUID `gorm:"column:id"`
			MetricKey string    `gorm:"column:metric_key"`
			GoalValue int       `gorm:"column:goal_value"`
		}
		if err := tx.WithContext(ctx).Raw(`
			SELECT id, metric_key, goal_value
			FROM achievements
			WHERE is_active = true
		`).Scan(&achievements).Error; err != nil {
			return err
		}

		if len(achievements) == 0 {
			return nil
		}

		totals, err := r.metricTotalsTx(ctx, tx, userID)
		if err != nil {
			return err
		}

		existingUnlockedAt, err := r.userAchievementUnlockedAtTx(ctx, tx, userID)
		if err != nil {
			return err
		}

		now := time.Now()
		for _, item := range achievements {
			progressValue := progressForMetric(totals, item.MetricKey)
			unlocked := progressValue >= item.GoalValue
			unlockedAt := existingUnlockedAt[item.ID]
			if unlocked && unlockedAt == nil {
				unlockedAt = &now
			}

			if err = tx.WithContext(ctx).Exec(`
				INSERT INTO user_achievements (
					id,
					user_id,
					achievement_id,
					progress_value,
					unlocked,
					unlocked_at,
					created_at,
					updated_at
				)
				VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
				ON CONFLICT (user_id, achievement_id)
				DO UPDATE SET
					progress_value = EXCLUDED.progress_value,
					unlocked = EXCLUDED.unlocked,
					unlocked_at = CASE
						WHEN user_achievements.unlocked_at IS NOT NULL THEN user_achievements.unlocked_at
						ELSE EXCLUDED.unlocked_at
					END,
					updated_at = NOW()
			`, uuid.New(), userID, item.ID, progressValue, unlocked, unlockedAt).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Repo) metricTotalsTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (metricTotals, error) {
	var totals metricTotals
	err := tx.WithContext(ctx).Raw(`
		SELECT
			(
				SELECT COUNT(DISTINCT lesson_id)::int
				FROM user_course_points
				WHERE user_id = ?
			) AS completed_lessons,
			(
				SELECT COUNT(DISTINCT quiz_id)::int
				FROM lesson_quiz_attempts
				WHERE user_id = ?
				  AND is_correct = true
			) AS passed_quizzes,
			(
				SELECT COALESCE((
					SELECT max_streak
					FROM users
					WHERE id = ?
				), 0)::int
			) AS max_streak,
			(
				SELECT COALESCE(SUM(xp), 0)::int
				FROM user_course_points
				WHERE user_id = ?
			) AS xp_total,
			(
				SELECT COUNT(*)::int
				FROM course_subscription
				WHERE user_id = ?
				  AND started_at IS NOT NULL
			) AS started_courses,
			(
				SELECT COUNT(*)::int
				FROM course_subscription
				WHERE user_id = ?
				  AND completed_at IS NOT NULL
			) AS completed_courses,
			(
				SELECT COUNT(*)::int
				FROM course_subscription
				WHERE user_id = ?
				  AND last_activity_at IS NOT NULL
			) AS active_courses,
			(
				SELECT COUNT(*)::int
				FROM (
					SELECT cm.id
					FROM course_modules cm
					INNER JOIN course_lessons cl ON cl.module_id = cm.id
					INNER JOIN user_course_points up ON up.lesson_id = cl.id AND up.user_id = ?
					GROUP BY cm.id
					HAVING COUNT(DISTINCT up.lesson_id) = COUNT(DISTINCT cl.id)
				) completed_modules
			) AS completed_modules
	`, userID, userID, userID, userID, userID, userID, userID, userID).Scan(&totals).Error
	return totals, err
}

func (r *Repo) userAchievementUnlockedAtTx(ctx context.Context, tx *gorm.DB, userID uuid.UUID) (map[uuid.UUID]*time.Time, error) {
	var rows []userAchievementRow
	if err := tx.WithContext(ctx).Raw(`
		SELECT achievement_id, unlocked_at
		FROM user_achievements
		WHERE user_id = ?
	`, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	items := make(map[uuid.UUID]*time.Time, len(rows))
	for _, row := range rows {
		items[row.AchievementID] = row.UnlockedAt
	}

	return items, nil
}

func progressForMetric(totals metricTotals, metricKey string) int {
	switch metricKey {
	case "completed_lessons":
		return totals.CompletedLessons
	case "passed_quizzes":
		return totals.PassedQuizzes
	case "max_streak":
		return totals.MaxStreak
	case "xp_total":
		return totals.XPTotal
	case "started_courses":
		return totals.StartedCourses
	case "completed_courses":
		return totals.CompletedCourses
	case "active_courses":
		return totals.ActiveCourses
	case "completed_modules":
		return totals.CompletedModules
	default:
		return 0
	}
}
