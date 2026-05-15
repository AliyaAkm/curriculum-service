package lesson

import "curriculum-service/internal/http/middleware"

type Handler struct {
	client      client
	localClient localClient
	videoStore  videoStorage
	jwtMgr      *middleware.Manager
}

func NewHandler(client client, localClient localClient, videoStore videoStorage, jwtMgr *middleware.Manager) *Handler {
	return &Handler{
		client:      client,
		localClient: localClient,
		videoStore:  videoStore,
		jwtMgr:      jwtMgr,
	}
}
