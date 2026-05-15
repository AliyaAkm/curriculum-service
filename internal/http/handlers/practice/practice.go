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

func (h *Handler) CreatePractice(c *gin.Context) {
	if !h.requireEditor(c) {
		return
	}

	var req practicedto.PracticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.client.CreatePractice(c.Request.Context(), convertPracticeRequest(req))
	if err != nil {
		writePracticeError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertPractice(result))
}

func (h *Handler) GetPracticeByLessonID(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Query("lesson_id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	userID, hasFullAccess, ok := h.currentUser(c)
	if !ok {
		return
	}

	result, err := h.client.GetPracticeByLessonIDForUser(c.Request.Context(), userID, lessonID, hasFullAccess)
	if err != nil {
		writePracticeError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertPractice(result))
}

func (h *Handler) GetPracticeByID(c *gin.Context) {
	practiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid practice id")
		return
	}

	userID, hasFullAccess, ok := h.currentUser(c)
	if !ok {
		return
	}

	result, err := h.client.GetPracticeByIDForUser(c.Request.Context(), userID, practiceID, hasFullAccess)
	if err != nil {
		writePracticeError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertPractice(result))
}

func (h *Handler) UpdatePractice(c *gin.Context) {
	if !h.requireEditor(c) {
		return
	}

	practiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid practice id")
		return
	}

	var req practicedto.PracticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.client.UpdatePractice(c.Request.Context(), practiceID, convertPracticeRequest(req))
	if err != nil {
		writePracticeError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertPractice(result))
}

func (h *Handler) DeletePractice(c *gin.Context) {
	if !h.requireEditor(c) {
		return
	}

	practiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid practice id")
		return
	}

	if err := h.client.DeletePractice(c.Request.Context(), practiceID); err != nil {
		writePracticeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) currentUser(c *gin.Context) (uuid.UUID, bool, bool) {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return uuid.Nil, false, false
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return uuid.Nil, false, false
	}

	return userID, middleware.ClaimsHasRole(claims, middleware.RoleAdmin), true
}

func (h *Handler) requireEditor(c *gin.Context) bool {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return false
	}
	if middleware.ClaimsHasRole(claims, middleware.RoleAdmin) || middleware.ClaimsHasRole(claims, middleware.RoleTeacher) {
		return true
	}

	respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
	return false
}

func writePracticeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request")
	case errors.Is(err, domain.ErrForbidden):
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
	case errors.Is(err, domain.ErrPracticeNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", domain.ErrPracticeNotFound.Error())
	default:
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}

func convertPracticeRequest(req practicedto.PracticeRequest) *practicedomain.Practice {
	return &practicedomain.Practice{
		LessonID:         req.LessonID,
		Position:         req.Position,
		Title:            convertLocaleRequest(req.Title),
		Summary:          convertLocaleRequest(req.Summary),
		Brief:            convertLocaleRequest(req.Brief),
		StarterCode:      req.StarterCode,
		SuccessCriteria:  convertLocaleRequests(req.SuccessCriteria),
		KnowledgeChecks:  convertLocaleRequests(req.KnowledgeChecks),
		PromptSuggestion: convertLocaleRequest(req.PromptSuggestion),
		XPReward:         req.XPReward,
	}
}

func convertPractice(resp *practicedomain.Practice) practicedto.Practice {
	if resp == nil {
		return practicedto.Practice{}
	}

	return practicedto.Practice{
		ID:               resp.ID,
		LessonID:         resp.LessonID,
		Position:         resp.Position,
		Title:            convertLocale(resp.Title),
		Summary:          convertLocale(resp.Summary),
		Brief:            convertLocale(resp.Brief),
		StarterCode:      resp.StarterCode,
		SuccessCriteria:  convertLocales(resp.SuccessCriteria),
		KnowledgeChecks:  convertLocales(resp.KnowledgeChecks),
		PromptSuggestion: convertLocale(resp.PromptSuggestion),
		XPReward:         resp.XPReward,
		CreatedAt:        resp.CreatedAt,
		UpdatedAt:        resp.UpdatedAt,
	}
}

func convertLocaleRequest(src practicedto.Locale) practicedomain.Locale {
	return practicedomain.Locale{
		EN: src.EN,
		RU: src.RU,
		KK: src.KK,
	}
}

func convertLocale(src practicedomain.Locale) practicedto.Locale {
	return practicedto.Locale{
		EN: src.EN,
		RU: src.RU,
		KK: src.KK,
	}
}

func convertLocaleRequests(src []practicedto.Locale) []practicedomain.Locale {
	result := make([]practicedomain.Locale, len(src))
	for i := range src {
		result[i] = convertLocaleRequest(src[i])
	}
	return result
}

func convertLocales(src []practicedomain.Locale) []practicedto.Locale {
	result := make([]practicedto.Locale, len(src))
	for i := range src {
		result[i] = convertLocale(src[i])
	}
	return result
}
