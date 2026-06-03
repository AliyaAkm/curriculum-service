package certificate

type UseCase struct {
	repo    Repository
	storage Storage
}

func NewUseCase(repo Repository, storage Storage) *UseCase {
	return &UseCase{
		repo:    repo,
		storage: storage,
	}
}
