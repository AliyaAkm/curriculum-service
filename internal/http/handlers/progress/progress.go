package progress

import (
	"curriculum-service/internal/domain"
	progressdomain "curriculum-service/internal/domain/progress"
	progressdto "curriculum-service/internal/http/dto/progress"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CompleteLesson(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}

	lessonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	result, err := h.client.CompleteLesson(c.Request.Context(), userID, lessonID)
	if err != nil {
		writeProgressError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertCourseProgress(result))
}

func (h *Handler) GetCourseProgress(c *gin.Context) {
	claims, ok := h.claims(c)
	if !ok {
		return
	}
	userID, ok := userIDFromClaims(c, claims)
	if !ok {
		return
	}

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid course id")
		return
	}

	requestedUserID := c.Param("user_id")
	if requestedUserID == "" {
		requestedUserID = c.Query("user_id")
	}

	if requestedUserID != "" {
		if !middleware.ClaimsHasRole(claims, middleware.RoleTeacher) && !middleware.ClaimsHasRole(claims, middleware.RoleAdmin) {
			respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
			return
		}

		userID, err = uuid.Parse(requestedUserID)
		if err != nil {
			respond.JSON(c, http.StatusBadRequest, "invalid user id")
			return
		}
	}

	result, err := h.client.GetCourseProgress(c.Request.Context(), userID, courseID)
	if err != nil {
		writeProgressError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertCourseProgress(result))
}

func (h *Handler) ListCourseProgress(c *gin.Context) {
	claims, ok := h.claims(c)
	if !ok {
		return
	}
	userID, ok := userIDFromClaims(c, claims)
	if !ok {
		return
	}

	requestedUserID := c.Param("user_id")
	if requestedUserID != "" {
		if !middleware.ClaimsHasRole(claims, middleware.RoleAdmin) {
			respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
			return
		}

		parsedUserID, err := uuid.Parse(requestedUserID)
		if err != nil {
			respond.JSON(c, http.StatusBadRequest, "invalid user id")
			return
		}
		userID = parsedUserID
	}

	items, err := h.client.ListCourseProgress(c.Request.Context(), userID)
	if err != nil {
		writeProgressError(c, err)
		return
	}

	response := make([]progressdto.CourseProgress, len(items))
	for i := range items {
		response[i] = convertCourseProgress(&items[i])
	}

	respond.JSON(c, http.StatusOK, response)
}

func (h *Handler) userID(c *gin.Context) (uuid.UUID, bool) {
	claims, ok := h.claims(c)
	if !ok {
		return uuid.Nil, false
	}

	return userIDFromClaims(c, claims)
}

func (h *Handler) claims(c *gin.Context) (*middleware.Claims, bool) {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return nil, false
	}

	return claims, true
}

func userIDFromClaims(c *gin.Context, claims *middleware.Claims) (uuid.UUID, bool) {
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return uuid.Nil, false
	}

	return userID, true
}

func writeProgressError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request")
	case errors.Is(err, domain.ErrLessonNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", domain.ErrLessonNotFound.Error())
	case errors.Is(err, domain.ErrCourseNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", domain.ErrCourseNotFound.Error())
	case errors.Is(err, domain.ErrCourseProgressNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", domain.ErrCourseProgressNotFound.Error())
	case errors.Is(err, domain.ErrForbidden):
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
	default:
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}

func convertCourseProgress(src *progressdomain.CourseProgress) progressdto.CourseProgress {
	if src == nil {
		return progressdto.CourseProgress{}
	}

	modules := make([]progressdto.ModuleProgress, len(src.Modules))
	for i := range src.Modules {
		modules[i] = progressdto.ModuleProgress{
			ModuleID:         src.Modules[i].ModuleID,
			Position:         src.Modules[i].Position,
			IsOpen:           src.Modules[i].IsOpen,
			TotalLessons:     src.Modules[i].TotalLessons,
			CompletedLessons: src.Modules[i].CompletedLessons,
		}
	}

	return progressdto.CourseProgress{
		CourseID:           src.CourseID,
		UserID:             src.UserID,
		StartedAt:          src.StartedAt,
		LastActivityAt:     src.LastActivityAt,
		CompletedAt:        src.CompletedAt,
		CurrentLessonID:    src.CurrentLessonID,
		TotalLessons:       src.TotalLessons,
		CompletedLessons:   src.CompletedLessons,
		ProgressPercent:    src.ProgressPercent,
		CompletedLessonIDs: src.CompletedLessonIDs,
		PassedQuizIDs:      src.PassedQuizIDs,
		Modules:            modules,
	}
}
