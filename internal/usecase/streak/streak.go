package streak

import (
	"context"
	achievementdomain "curriculum-service/internal/domain/achievement"
	"curriculum-service/internal/domain/streak"
	"time"

	"github.com/google/uuid"
)

func (u *UseCase) GetStreak(ctx context.Context, userID uuid.UUID) (*streak.DailyStreak, error) {
	entity, err := u.repo.GetStreak(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if entity == nil {
		entity = &streak.DailyStreak{
			ID:        uuid.New(),
			UserID:    userID,
			Streak:    1,
			LastLogin: now,
		}
		if err := u.repo.CreateStreak(ctx, entity); err != nil {
			return nil, err
		}
		if err := u.repo.UpdateUserMaxStreak(ctx, userID, entity.Streak); err != nil {
			return nil, err
		}
		u.sendDailyMissionCompleted(ctx, userID, entity.Streak)
		return entity, nil
	}

	days := int(now.Sub(entity.LastLogin).Hours() / 24)
	if days == 0 {
		if err := u.repo.UpdateUserMaxStreak(ctx, userID, entity.Streak); err != nil {
			return nil, err
		}
		return entity, nil
	}
	if days == 1 {
		entity.Streak++
	} else {
		entity.Streak = 1
	}
	entity.LastLogin = now

	if err := u.repo.UpdateStreak(ctx, entity); err != nil {
		return nil, err
	}
	if err := u.repo.UpdateUserMaxStreak(ctx, userID, entity.Streak); err != nil {
		return nil, err
	}
	u.sendDailyMissionCompleted(ctx, userID, entity.Streak)
	return entity, nil
}

func (u *UseCase) sendDailyMissionCompleted(ctx context.Context, userID uuid.UUID, streakValue int64) {
	if u.notification == nil {
		return
	}

	_ = u.notification.SendEvent(ctx, userID, "daily_mission_completed", map[string]any{
		"streak": streakValue,
	})
	u.notifyUnlockedAchievements(ctx, userID)
}

func (u *UseCase) notifyUnlockedAchievements(ctx context.Context, userID uuid.UUID) {
	if u.notification == nil || u.achievements == nil {
		return
	}

	items, err := u.achievements.SyncUnlockedAchievements(ctx, userID)
	if err != nil {
		return
	}
	for _, item := range items {
		_ = u.notification.SendEvent(ctx, userID, "achievement_unlocked", map[string]any{
			"achievementId":    item.ID.String(),
			"achievementCode":  item.Code,
			"achievementTitle": achievementTitle(item),
		})
	}
}

func achievementTitle(value achievementdomain.Achievement) string {
	switch {
	case value.Title.RU != "":
		return value.Title.RU
	case value.Title.EN != "":
		return value.Title.EN
	case value.Title.KK != "":
		return value.Title.KK
	default:
		return "новое достижение"
	}
}
