package progress

type UseCase struct {
	repo         Repository
	notification NotificationSender
}

func New(repo Repository, notification ...NotificationSender) *UseCase {
	var sender NotificationSender
	if len(notification) > 0 {
		sender = notification[0]
	}
	return &UseCase{
		repo:         repo,
		notification: sender,
	}
}
