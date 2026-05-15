package achievement

type UseCase struct {
	repo client
}

func New(repo client) *UseCase {
	return &UseCase{repo: repo}
}
