package streak

import (
	"curriculum-service/internal/domain"
	"curriculum-service/internal/domain/streak"
	dtostreak "curriculum-service/internal/http/dto/streak"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetStreak(c *gin.Context) {
	userID := middleware.GetUserID(h.jwtMgr, c)
	if userID == nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return
	}
	result, err := h.client.GetStreak(c.Request.Context(), *userID)
	if err != nil {
		writeStreakError(c, err)
		return
	}
	respond.JSON(c, http.StatusOK, convertStreak(result))
}

func convertStreak(resp *streak.DailyStreak) dtostreak.DailyStreakResponse {
	if resp == nil {
		return dtostreak.DailyStreakResponse{}
	}
	return dtostreak.DailyStreakResponse{
		Streak: resp.Streak,
	}
}

func writeStreakError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request query")
	default:
		c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}
