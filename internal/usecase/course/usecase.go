package course

type UseCase struct {
	repo         Repository
	reviewRepo   ReviewRepository
	moduleRepo   ModuleRepository
	notification NotificationSender
}

func New(repo Repository, reviewRepo ReviewRepository, moduleRepo ModuleRepository, notification ...NotificationSender) *UseCase {
	var sender NotificationSender
	if len(notification) > 0 {
		sender = notification[0]
	}
	return &UseCase{
		repo:         repo,
		reviewRepo:   reviewRepo,
		moduleRepo:   moduleRepo,
		notification: sender,
	}
}
