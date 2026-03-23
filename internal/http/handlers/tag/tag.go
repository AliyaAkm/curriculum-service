package tag

import (
	"curriculum-service/internal/domain"
	domaintag "curriculum-service/internal/domain/tag"
	dtotag "curriculum-service/internal/http/dto/tag"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListCourseTags(c *gin.Context) {
	result, err := h.client.GetAllTags(c.Request.Context())
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertTags(result))
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

func convertTags(resp []domaintag.Tag) []dtotag.Tag {
	tags := make([]dtotag.Tag, len(resp))

	for i := range resp {
		tags[i] = dtotag.Tag{
			ID:   resp[i].ID,
			Name: resp[i].Name,
			Code: resp[i].Code,
		}
	}

	return tags
}
