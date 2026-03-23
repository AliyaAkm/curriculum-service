package level

import (
	"curriculum-service/internal/domain"
	domainlevel "curriculum-service/internal/domain/level"
	dtolevels "curriculum-service/internal/http/dto/level"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListCourseLevels(c *gin.Context) {
	result, err := h.client.GetAllLevels(c.Request.Context())
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertLevels(result))
}

func writeCatalogError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request query")
	default:
		c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}

func convertLevels(resp []domainlevel.Level) []dtolevels.Level {
	levels := make([]dtolevels.Level, len(resp))

	for i := range resp {
		levels[i] = dtolevels.Level{
			ID:   resp[i].ID,
			Name: resp[i].Name,
			Code: resp[i].Code,
		}
	}

	return levels
}
