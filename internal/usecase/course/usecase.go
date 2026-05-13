package course

type UseCase struct {
	repo       Repository
	reviewRepo ReviewRepository
	moduleRepo ModuleRepository
}

func New(repo Repository, reviewRepo ReviewRepository, moduleRepo ModuleRepository) *UseCase {
	return &UseCase{
		repo:       repo,
		reviewRepo: reviewRepo,
		moduleRepo: moduleRepo,
	}
}
