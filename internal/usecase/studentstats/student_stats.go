package studentstats

import (
	"context"
	studentstatsdomain "curriculum-service/internal/domain/studentstats"
	"time"

	"github.com/google/uuid"
)

func (u *UseCase) GetStatistics(ctx context.Context, userID uuid.UUID) (*studentstatsdomain.Statistics, error) {
	stats, err := u.repo.GetStatistics(ctx, userID)
	if err != nil {
		return nil, err
	}

	if u.ai != nil {
		aiStats, aiActivity, err := u.ai.GetUserAnalytics(ctx, userID)
		if err == nil && aiStats != nil {
			stats.AI = *aiStats
			mergeAIActivity(stats.Activity, aiActivity)
			if newerTime(stats.Summary.LastActivityAt, aiStats.LastActivityAt) {
				stats.Summary.LastActivityAt = aiStats.LastActivityAt
			}
		}
	}

	stats.Recommendations = buildRecommendations(stats)
	return stats, nil
}

func mergeAIActivity(target []studentstatsdomain.ActivityDay, source []studentstatsdomain.ActivityDay) {
	byDate := make(map[string]int, len(target))
	for i := range target {
		byDate[target[i].Date] = i
	}
	for _, item := range source {
		if idx, ok := byDate[item.Date]; ok {
			target[idx].AIRequests += item.AIRequests
		}
	}
}

func buildRecommendations(stats *studentstatsdomain.Statistics) []studentstatsdomain.Recommendation {
	items := make([]studentstatsdomain.Recommendation, 0, 3)
	if stats.Summary.ActiveCourses == 0 && stats.Summary.StartedCourses > 0 {
		items = append(items, studentstatsdomain.Recommendation{
			Type:    "course",
			Title:   "Вернуться к курсу",
			Message: "Есть начатые курсы без свежей активности.",
		})
	}
	if stats.Quiz.AttemptedQuizzes > 0 && stats.Quiz.AccuracyPercent < 70 {
		items = append(items, studentstatsdomain.Recommendation{
			Type:    "quiz",
			Title:   "Повторить теорию",
			Message: "Точность квизов ниже 70%, лучше закрепить сложные уроки.",
		})
	}
	if stats.Practice.Submissions > 0 && stats.Practice.PassRatePercent < 60 {
		items = append(items, studentstatsdomain.Recommendation{
			Type:    "practice",
			Title:   "Разобрать практику",
			Message: "Практические задания часто не проходят с первого раза.",
		})
	}
	return items
}

func newerTime(current, candidate *time.Time) bool {
	if candidate == nil {
		return false
	}
	if current == nil {
		return true
	}
	return current.Before(*candidate)
}
