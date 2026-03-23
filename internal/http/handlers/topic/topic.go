package topic

import (
	"curriculum-service/internal/domain"
	domaintopic "curriculum-service/internal/domain/topic"
	dtotopic "curriculum-service/internal/http/dto/topic"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListCourseTopics(c *gin.Context) {
	result, err := h.client.GetAllTopics(c.Request.Context())
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertTopics(result))
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

func convertTopics(resp []domaintopic.Topic) []dtotopic.Topic {
	topics := make([]dtotopic.Topic, len(resp))

	for i := range resp {
		topics[i] = dtotopic.Topic{
			ID:   resp[i].ID,
			Name: resp[i].Name,
			Code: resp[i].Code,
		}
	}

	return topics
}
