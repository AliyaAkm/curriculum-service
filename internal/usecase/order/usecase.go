package order

type UseCase struct {
	repo       Repository
	statusRepo StatusRepository
	priceRepo  PriceRepository
}

func New(repo Repository, statusRepo StatusRepository, priceRepo PriceRepository) *UseCase {
	return &UseCase{
		repo:       repo,
		statusRepo: statusRepo,
		priceRepo:  priceRepo,
	}
}
