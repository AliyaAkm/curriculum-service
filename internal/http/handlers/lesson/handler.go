package lesson

type Handler struct {
	client      client
	localClient localClient
}

func NewHandler(client client, localClient localClient) *Handler {
	return &Handler{
		client:      client,
		localClient: localClient,
	}
}
