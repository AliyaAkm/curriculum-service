package streak

type UseCase struct {
	repo         Repository
	notification NotificationSender
	achievements AchievementSyncer
}

func New(repo Repository, notification NotificationSender, achievements ...AchievementSyncer) *UseCase {
	var achievementSyncer AchievementSyncer
	if len(achievements) > 0 {
		achievementSyncer = achievements[0]
	}
	return &UseCase{
		repo:         repo,
		notification: notification,
		achievements: achievementSyncer,
	}
}
