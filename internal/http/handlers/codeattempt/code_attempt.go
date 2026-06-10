package codeattempt

import (
	"curriculum-service/internal/domain"
	codeattemptdomain "curriculum-service/internal/domain/codeattempt"
	codeattemptdto "curriculum-service/internal/http/dto/codeattempt"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Run(c *gin.Context) {
	userID, ok := h.userID(c)
	if !ok {
		return
	}

	var req codeattemptdto.RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	var courseID *uuid.UUID
	if req.CourseID != uuid.Nil {
		courseID = &req.CourseID
	}
	var lessonID *uuid.UUID
	if req.LessonID != uuid.Nil {
		lessonID = &req.LessonID
	}

	result, err := h.client.Run(c.Request.Context(), codeattemptdomain.RunRequest{
		UserID:     userID,
		CourseID:   courseID,
		LessonID:   lessonID,
		PracticeID: c.Param("id"),
		RunType:    req.RunType,
		Language:   req.Language,
		Code:       req.Code,
	})
	if err != nil {
		writeError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, codeattemptdto.RunResponse{
		AttemptID:  result.AttemptID,
		Output:     result.Output,
		Error:      result.Error,
		Passed:     result.Passed,
		ErrorType:  result.ErrorType,
		DurationMS: result.DurationMS,
		XPAwarded:  result.XPAwarded,
	})
}

func (h *Handler) userID(c *gin.Context) (uuid.UUID, bool) {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return uuid.Nil, false
	}

	return userID, true
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request")
	case errors.Is(err, domain.ErrPracticeAutoSubmitNotAllowed):
		respond.Error(c, http.StatusBadRequest, "validation", domain.ErrPracticeAutoSubmitNotAllowed.Error())
	case errors.Is(err, domain.ErrPracticePrerequisitesNotMet):
		respond.Error(c, http.StatusConflict, "practice_prerequisites_not_met", domain.ErrPracticePrerequisitesNotMet.Error())
	case errors.Is(err, domain.ErrPracticeNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", domain.ErrPracticeNotFound.Error())
	default:
		_ = c.Error(err)
		respond.Error(c, http.StatusBadGateway, "code_runner_error", err.Error())
	}
}
