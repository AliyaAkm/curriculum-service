package certificate

import "curriculum-service/internal/http/middleware"

type Handler struct {
	client client
	jwtMgr *middleware.Manager
}

func NewHandler(client client, jwtMgr *middleware.Manager) *Handler {
	return &Handler{
		client: client,
		jwtMgr: jwtMgr,
	}
}
