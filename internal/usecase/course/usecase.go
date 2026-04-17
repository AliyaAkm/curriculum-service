package course

type UseCase struct {
	repo       Repository
	reviewRepo ReviewRepository
}

func New(repo Repository, reviewRepo ReviewRepository) *UseCase {
	return &UseCase{
		repo:       repo,
		reviewRepo: reviewRepo,
	}
}
