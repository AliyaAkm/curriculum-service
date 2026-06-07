package studentstats

type UseCase struct {
	repo Repository
	ai   AIAnalyticsClient
}

func New(repo Repository, ai AIAnalyticsClient) *UseCase {
	return &UseCase{
		repo: repo,
		ai:   ai,
	}
}
