package studentstats

import (
	"curriculum-service/internal/domain"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) GetStatistics(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}

	stats, err := h.client.GetStatistics(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
		return
	}

	respond.JSON(c, http.StatusOK, stats)
}

func (h *Handler) userID(c *gin.Context) (uuid.UUID, bool) {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return uuid.Nil, false
	}

	return userID, true
}
