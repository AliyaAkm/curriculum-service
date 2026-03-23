package course

type Handler struct {
	client client
}

func New(client client) *Handler {
	return &Handler{
		client: client,
	}
}
