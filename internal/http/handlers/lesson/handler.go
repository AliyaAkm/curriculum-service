package lesson

import "curriculum-service/internal/http/middleware"

type Handler struct {
	client      client
	localClient localClient
	storage     objectStorage
	jwtMgr      *middleware.Manager
}

func NewHandler(client client, localClient localClient, storage objectStorage, jwtMgr *middleware.Manager) *Handler {
	return &Handler{
		client:      client,
		localClient: localClient,
		storage:     storage,
		jwtMgr:      jwtMgr,
	}
}
