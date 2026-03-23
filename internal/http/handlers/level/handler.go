package level

type Handler struct {
	client client
}

func NewHandler(client client) *Handler {
	return &Handler{
		client: client,
	}
}
