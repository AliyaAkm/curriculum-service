package locale

import (
	"curriculum-service/internal/domain"
	domainlocale "curriculum-service/internal/domain/locale"
	dtolocale "curriculum-service/internal/http/dto/locale"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListCourseLocales(c *gin.Context) {
	result, err := h.client.GetAllLocales(c.Request.Context())
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertLocales(result))
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

func convertLocales(resp []domainlocale.Locale) []dtolocale.Locale {
	locales := make([]dtolocale.Locale, len(resp))

	for i := range resp {
		locales[i] = dtolocale.Locale{
			ID:   resp[i].ID,
			Name: resp[i].Name,
			Code: resp[i].Code,
		}
	}

	return locales
}
