package achievement

import (
	"curriculum-service/internal/domain"
	achievementdomain "curriculum-service/internal/domain/achievement"
	achievementdto "curriculum-service/internal/http/dto/achievement"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListAchievements(c *gin.Context) {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return
	}

	requestedUserID := c.Param("user_id")
	if requestedUserID != "" {
		if !middleware.ClaimsHasRole(claims, middleware.RoleAdmin) {
			respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
			return
		}

		userID, err = uuid.Parse(requestedUserID)
		if err != nil {
			respond.JSON(c, http.StatusBadRequest, "invalid user id")
			return
		}
	}

	items, err := h.client.ListAchievements(c.Request.Context(), userID)
	if err != nil {
		writeAchievementError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertAchievements(items))
}

func writeAchievementError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrForbidden):
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request")
	default:
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}

func convertAchievements(items []achievementdomain.Achievement) []achievementdto.Achievement {
	response := make([]achievementdto.Achievement, len(items))
	for i := range items {
		response[i] = achievementdto.Achievement{
			ID: items[i].Code,
			Title: achievementdto.LocalizedText{
				EN: items[i].Title.EN,
				RU: items[i].Title.RU,
				KK: items[i].Title.KK,
			},
			Description: achievementdto.LocalizedText{
				EN: items[i].Description.EN,
				RU: items[i].Description.RU,
				KK: items[i].Description.KK,
			},
			IconKey:    items[i].IconKey,
			Goal:       items[i].Goal,
			Progress:   items[i].Progress,
			Unlocked:   items[i].Unlocked,
			UnlockedAt: items[i].UnlockedAt,
		}
	}

	return response
}
