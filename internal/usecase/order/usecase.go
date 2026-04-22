package order

type UseCase struct {
	repo       Repository
	statusRepo StatusRepository
}

func New(repo Repository, statusRepo StatusRepository) *UseCase {
	return &UseCase{
		repo:       repo,
		statusRepo: statusRepo,
	}
}
