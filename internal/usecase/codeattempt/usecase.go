package codeattempt

type UseCase struct {
	repo      Repository
	runner    Runner
	practices PracticeProvider
}

func New(repo Repository, runner Runner, practices PracticeProvider) *UseCase {
	return &UseCase{
		repo:      repo,
		runner:    runner,
		practices: practices,
	}
}
