package handlers

import (
	"curriculum-service/internal/domain/category"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"curriculum-service/internal/domain"
	"curriculum-service/internal/http/respond"
	"curriculum-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	uc *usecase.Catalog
}

func NewCatalogHandler(uc *usecase.Catalog) *CatalogHandler {
	return &CatalogHandler{uc: uc}
}

func (h *CatalogHandler) ListCourses(c *gin.Context) {
	if hasTextSearch(c) {
		respond.Error(c, http.StatusBadRequest, "validation", "use /course/search for text search")
		return
	}
	if hasStructuredFilters(c) {
		respond.Error(c, http.StatusBadRequest, "validation", "use /course/filter for structured filters")
		return
	}

	filter := category.CourseSearchFilter{
		Locale:   category.NormalizeLocale(c.DefaultQuery("locale", "ru")),
		Page:     parsePositiveInt(c.Query("page"), 1),
		PageSize: parsePositiveInt(c.DefaultQuery("page_size", "12"), 12),
	}

	result, err := h.uc.SearchCourses(c.Request.Context(), filter)
	if err != nil {
		fmt.Println(err)
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, toSearchCoursesResponse(result))
}

func (h *CatalogHandler) SearchCourses(c *gin.Context) {
	if hasStructuredFilters(c) {
		respond.Error(c, http.StatusBadRequest, "validation", "use /course for structured filters")
		return
	}

	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		respond.Error(c, http.StatusBadRequest, "validation", "q is required for /course/search")
		return
	}

	filter := category.CourseSearchFilter{
		Query:    query,
		Locale:   category.NormalizeLocale(c.DefaultQuery("locale", "ru")),
		Page:     parsePositiveInt(c.Query("page"), 1),
		PageSize: parsePositiveInt(c.DefaultQuery("page_size", "12"), 12),
	}

	result, err := h.uc.SearchCourses(c.Request.Context(), filter)
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, toSearchCoursesResponse(result))
}

func collectQueryValues(c *gin.Context, keys ...string) []string {
	query := c.Request.URL.Query()
	result := make([]string, 0)
	for _, key := range keys {
		for _, raw := range query[key] {
			for _, item := range strings.Split(raw, ",") {
				trimmed := strings.TrimSpace(item)
				if trimmed == "" {
					continue
				}
				result = append(result, trimmed)
			}
		}
	}
	return result
}

func hasTextSearch(c *gin.Context) bool {
	return strings.TrimSpace(c.Query("q")) != ""
}

func hasStructuredFilters(c *gin.Context) bool {
	query := c.Request.URL.Query()
	keys := []string{
		"topic",
		"topics",
		"level",
		"levels",
		"min_rating",
		"duration",
		"durations",
		"with_certificate",
	}

	for _, key := range keys {
		for _, value := range query[key] {
			if strings.TrimSpace(value) != "" {
				return true
			}
		}
	}

	return false
}

func parseMinRating(raw string) (float64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, domain.ErrValidation
	}

	return value, nil
}

func parseWithCertificate(raw string) (*bool, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, domain.ErrValidation
	}

	if !value {
		return nil, nil
	}

	return &value, nil
}

func parsePositiveInt(raw string, fallback int) int {
	value, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
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
