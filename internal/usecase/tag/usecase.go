package tag

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}
