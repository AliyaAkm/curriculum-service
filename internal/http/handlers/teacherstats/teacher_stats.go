package teacherstats

import (
	"curriculum-service/internal/domain"
	teacherstatsdomain "curriculum-service/internal/domain/teacherstats"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) GetStatistics(c *gin.Context) {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return
	}
	if !middleware.ClaimsHasRole(claims, middleware.RoleTeacher) && !middleware.ClaimsHasRole(claims, middleware.RoleAdmin) {
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
		return
	}

	teacherID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return
	}

	courseID, ok := optionalUUIDQuery(c, "course_id")
	if !ok {
		return
	}

	stats, err := h.client.GetStatistics(c.Request.Context(), teacherstatsdomain.Filter{
		TeacherID:  teacherID,
		IsAdmin:    middleware.ClaimsHasRole(claims, middleware.RoleAdmin),
		CourseID:   courseID,
		PeriodDays: intQuery(c, "period_days", 30),
	})
	if err != nil {
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
		return
	}

	respond.JSON(c, http.StatusOK, stats)
}

func optionalUUIDQuery(c *gin.Context, key string) (*uuid.UUID, bool) {
	value := c.Query(key)
	if value == "" {
		return nil, true
	}
	parsed, err := uuid.Parse(value)
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "validation", "invalid "+key)
		return nil, false
	}
	return &parsed, true
}

func intQuery(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
