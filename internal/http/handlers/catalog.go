package handlers

import (
	"errors"
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
		respond.Error(c, http.StatusBadRequest, "validation", "use /courses/search for text search")
		return
	}
	if hasStructuredFilters(c) {
		respond.Error(c, http.StatusBadRequest, "validation", "use /courses/filter for structured filters")
		return
	}

	filter := domain.CourseSearchFilter{
		Locale:   domain.NormalizeLocale(c.DefaultQuery("locale", "ru")),
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

func (h *CatalogHandler) SearchCourses(c *gin.Context) {
	if hasStructuredFilters(c) {
		respond.Error(c, http.StatusBadRequest, "validation", "use /courses for structured filters")
		return
	}

	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		respond.Error(c, http.StatusBadRequest, "validation", "q is required for /courses/search")
		return
	}

	filter := domain.CourseSearchFilter{
		Query:    query,
		Locale:   domain.NormalizeLocale(c.DefaultQuery("locale", "ru")),
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

func (h *CatalogHandler) FilterCourses(c *gin.Context) {
	if hasTextSearch(c) {
		respond.Error(c, http.StatusBadRequest, "validation", "use /courses/search for text search")
		return
	}
	if !hasStructuredFilters(c) {
		respond.Error(c, http.StatusBadRequest, "validation", "at least one filter is required for /courses/filter")
		return
	}

	filter, err := buildStructuredFilter(c)
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	result, err := h.uc.SearchCourses(c.Request.Context(), filter)
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, toSearchCoursesResponse(result))
}

func (h *CatalogHandler) ListFilterOptions(c *gin.Context) {
	options, err := h.uc.GetFilterOptions(c.Request.Context(), domain.NormalizeLocale(c.DefaultQuery("locale", "ru")))
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, toFilterOptionsResponse(options))
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

func buildStructuredFilter(c *gin.Context) (domain.CourseSearchFilter, error) {
	levels, err := domain.ParseCourseLevels(collectQueryValues(c, "level", "levels"))
	if err != nil {
		return domain.CourseSearchFilter{}, err
	}

	durations, err := domain.ParseDurationCategories(collectQueryValues(c, "duration", "durations"))
	if err != nil {
		return domain.CourseSearchFilter{}, err
	}

	minRating, err := parseMinRating(c.Query("min_rating"))
	if err != nil {
		return domain.CourseSearchFilter{}, err
	}

	withCertificate, err := parseWithCertificate(c.Query("with_certificate"))
	if err != nil {
		return domain.CourseSearchFilter{}, err
	}

	return domain.CourseSearchFilter{
		Locale:          domain.NormalizeLocale(c.DefaultQuery("locale", "ru")),
		TopicSlugs:      collectQueryValues(c, "topic", "topics"),
		Levels:          levels,
		MinRating:       minRating,
		Durations:       durations,
		WithCertificate: withCertificate,
		Page:            parsePositiveInt(c.Query("page"), 1),
		PageSize:        parsePositiveInt(c.DefaultQuery("page_size", "12"), 12),
	}, nil
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
