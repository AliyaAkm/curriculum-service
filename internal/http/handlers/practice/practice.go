package practice

import (
	"curriculum-service/internal/domain"
	practicedomain "curriculum-service/internal/domain/practice"
	practicedto "curriculum-service/internal/http/dto/practice"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Create(c *gin.Context) {
	if !h.requireManager(c) {
		return
	}

	var req practicedto.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	result, err := h.client.Create(c.Request.Context(), taskFromRequest(req))
	if err != nil {
		writePracticeError(c, err)
		return
	}

	respond.JSON(c, http.StatusCreated, taskResponse(result, true))
}

func (h *Handler) Update(c *gin.Context) {
	if !h.requireManager(c) {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "validation", "invalid practice id")
		return
	}

	var req practicedto.TaskUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	result, err := h.client.Update(c.Request.Context(), id, taskUpdateFromRequest(req))
	if err != nil {
		writePracticeError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, taskResponse(result, true))
}

func (h *Handler) Delete(c *gin.Context) {
	if !h.requireManager(c) {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "validation", "invalid practice id")
		return
	}

	if err = h.client.Delete(c.Request.Context(), id); err != nil {
		writePracticeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetByID(c *gin.Context) {
	if !h.requireAuthenticated(c) {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "validation", "invalid practice id")
		return
	}

	result, err := h.client.GetByID(c.Request.Context(), id)
	if err != nil {
		writePracticeError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, taskResponse(result, h.canManage(c)))
}

func (h *Handler) ListByLesson(c *gin.Context) {
	if !h.requireAuthenticated(c) {
		return
	}

	lessonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "validation", "invalid lesson id")
		return
	}

	items, err := h.client.ListByLesson(c.Request.Context(), lessonID)
	if err != nil {
		writePracticeError(c, err)
		return
	}

	includeExpected := h.canManage(c)
	response := make([]practicedto.TaskResponse, len(items))
	for i := range items {
		response[i] = taskResponse(&items[i], includeExpected)
	}

	respond.JSON(c, http.StatusOK, response)
}

func (h *Handler) requireAuthenticated(c *gin.Context) bool {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return false
	}
	return true
}

func (h *Handler) requireManager(c *gin.Context) bool {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return false
	}
	if middleware.ClaimsHasRole(claims, middleware.RoleTeacher) {
		return true
	}

	respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
	return false
}

func (h *Handler) canManage(c *gin.Context) bool {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		return false
	}
	return middleware.ClaimsHasRole(claims, middleware.RoleTeacher) || middleware.ClaimsHasRole(claims, middleware.RoleAdmin)
}

func taskFromRequest(req practicedto.TaskRequest) practicedomain.Task {
	return practicedomain.Task{
		LessonID:       req.LessonID,
		Position:       req.Position,
		Title:          req.Title,
		Description:    req.Description,
		Language:       req.Language,
		StarterCode:    req.StarterCode,
		ExpectedOutput: req.ExpectedOutput,
		XPReward:       req.XPReward,
	}
}

func taskUpdateFromRequest(req practicedto.TaskUpdateRequest) practicedomain.TaskUpdate {
	return practicedomain.TaskUpdate{
		Title:          req.Title,
		Description:    req.Description,
		Language:       req.Language,
		StarterCode:    req.StarterCode,
		ExpectedOutput: req.ExpectedOutput,
		XPReward:       req.XPReward,
	}
}

func taskResponse(src *practicedomain.Task, includeExpected bool) practicedto.TaskResponse {
	if src == nil {
		return practicedto.TaskResponse{}
	}

	response := practicedto.TaskResponse{
		ID:          src.ID,
		LessonID:    src.LessonID,
		ModuleID:    src.ModuleID,
		CourseID:    src.CourseID,
		Position:    src.Position,
		Title:       src.Title,
		Description: src.Description,
		Language:    src.Language,
		StarterCode: src.StarterCode,
		XPReward:    src.XPReward,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
	if includeExpected {
		response.ExpectedOutput = src.ExpectedOutput
	}
	return response
}

func writePracticeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request")
	case errors.Is(err, domain.ErrPracticeNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", domain.ErrPracticeNotFound.Error())
	default:
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}
