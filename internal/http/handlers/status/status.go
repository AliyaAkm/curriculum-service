package status

import (
	"curriculum-service/internal/domain"
	domainstatus "curriculum-service/internal/domain/status"
	dtostatus "curriculum-service/internal/http/dto/status"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListCourseStatuses(c *gin.Context) {
	result, err := h.client.GetAllStatuses(c.Request.Context())
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertStatuses(result))
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

func convertStatuses(resp []domainstatus.Status) []dtostatus.Status {
	statuses := make([]dtostatus.Status, len(resp))

	for i := range resp {
		statuses[i] = dtostatus.Status{
			ID:   resp[i].ID,
			Name: resp[i].Name,
			Code: resp[i].Code,
		}
	}

	return statuses
}
