package durationcategory

import (
	"curriculum-service/internal/domain"
	domaindurationcategory "curriculum-service/internal/domain/durationcategory"
	dtodurationcategory "curriculum-service/internal/http/dto/durationcategory"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListCourseDurationCategories(c *gin.Context) {
	result, err := h.client.GetAllDurationCategories(c.Request.Context())
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertDurationCategory(result))
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

func convertDurationCategory(resp []domaindurationcategory.DurationCategory) []dtodurationcategory.DurationCategory {
	durationCategory := make([]dtodurationcategory.DurationCategory, len(resp))

	for i := range resp {
		durationCategory[i] = dtodurationcategory.DurationCategory{
			ID:   resp[i].ID,
			Name: resp[i].Name,
			Code: resp[i].Code,
		}
	}

	return durationCategory
}
