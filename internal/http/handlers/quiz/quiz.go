package quiz

import (
	"curriculum-service/internal/domain"
	quizdomain "curriculum-service/internal/domain/quiz"
	quizdto "curriculum-service/internal/http/dto/quiz"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateQuiz(c *gin.Context) {
	if !h.requireEditor(c) {
		return
	}

	var req quizdto.QuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.client.CreateQuiz(c.Request.Context(), convertQuizRequest(req))
	if err != nil {
		writeQuizError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertQuiz(result))
}

func (h *Handler) GetQuizzesByLessonID(c *gin.Context) {
	lessonID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid lesson id")
		return
	}

	userID, hasFullAccess, ok := h.currentUser(c)
	if !ok {
		return
	}

	result, err := h.client.GetQuizzesByLessonIDForUser(c.Request.Context(), userID, lessonID, hasFullAccess)
	if err != nil {
		writeQuizError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertQuizzes(result))
}

func (h *Handler) GetQuizByID(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid quiz id")
		return
	}

	userID, hasFullAccess, ok := h.currentUser(c)
	if !ok {
		return
	}

	result, err := h.client.GetQuizByIDForUser(c.Request.Context(), userID, quizID, hasFullAccess)
	if err != nil {
		writeQuizError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertQuiz(result))
}

func (h *Handler) SubmitAnswer(c *gin.Context) {
	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid quiz id")
		return
	}

	userID, hasFullAccess, ok := h.currentUser(c)
	if !ok {
		return
	}

	var req quizdto.AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.client.SubmitAnswer(c.Request.Context(), userID, quizID, req.SelectedAnswerIndex, hasFullAccess)
	if err != nil {
		writeQuizError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertAnswerResult(result))
}

func (h *Handler) UpdateQuiz(c *gin.Context) {
	if !h.requireEditor(c) {
		return
	}

	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid quiz id")
		return
	}

	var req quizdto.QuizUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.client.UpdateQuiz(c.Request.Context(), quizID, convertQuizUpdateRequest(req))
	if err != nil {
		writeQuizError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertQuiz(result))
}

func (h *Handler) DeleteQuiz(c *gin.Context) {
	if !h.requireEditor(c) {
		return
	}

	quizID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid quiz id")
		return
	}

	if err := h.client.DeleteQuiz(c.Request.Context(), quizID); err != nil {
		writeQuizError(c, err)
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

func writeQuizError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request")
	case errors.Is(err, domain.ErrForbidden):
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
	case errors.Is(err, domain.ErrQuizNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", domain.ErrQuizNotFound.Error())
	default:
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}

func convertQuizRequest(req quizdto.QuizRequest) *quizdomain.Quiz {
	options := make([]quizdomain.Option, len(req.Options))
	for i := range req.Options {
		options[i] = quizdomain.Option{
			Text: convertLocaleRequest(req.Options[i]),
		}
	}

	return &quizdomain.Quiz{
		LessonID:           req.LessonID,
		Position:           req.Position,
		Question:           convertLocaleRequest(req.Question),
		Options:            options,
		CorrectAnswerIndex: req.CorrectAnswerIndex,
		Explanation:        convertLocaleRequest(req.Explanation),
	}
}

func convertQuizUpdateRequest(req quizdto.QuizUpdateRequest) *quizdomain.Quiz {
	options := make([]quizdomain.Option, len(req.Options))
	for i := range req.Options {
		options[i] = quizdomain.Option{
			Text: convertLocaleRequest(req.Options[i]),
		}
	}

	return &quizdomain.Quiz{
		Position:           req.Position,
		Question:           convertLocaleRequest(req.Question),
		Options:            options,
		CorrectAnswerIndex: req.CorrectAnswerIndex,
		Explanation:        convertLocaleRequest(req.Explanation),
	}
}

func convertQuizzes(resp []quizdomain.Quiz) []quizdto.Quiz {
	quizzes := make([]quizdto.Quiz, len(resp))
	for i := range resp {
		quizzes[i] = convertQuiz(&resp[i])
	}
	return quizzes
}

func convertQuiz(resp *quizdomain.Quiz) quizdto.Quiz {
	if resp == nil {
		return quizdto.Quiz{}
	}

	options := make([]quizdto.Option, len(resp.Options))
	for i := range resp.Options {
		options[i] = quizdto.Option{
			ID:       resp.Options[i].ID,
			Position: resp.Options[i].Position,
			Text:     convertLocale(resp.Options[i].Text),
		}
	}

	var correctOptionID uuid.UUID
	if resp.CorrectAnswerIndex >= 0 && resp.CorrectAnswerIndex < len(resp.Options) {
		correctOptionID = resp.Options[resp.CorrectAnswerIndex].ID
	}

	return quizdto.Quiz{
		ID:                 resp.ID,
		LessonID:           resp.LessonID,
		Position:           resp.Position,
		Question:           convertLocale(resp.Question),
		Options:            options,
		CorrectAnswerIndex: resp.CorrectAnswerIndex,
		CorrectOptionID:    correctOptionID,
		Explanation:        convertLocale(resp.Explanation),
		CreatedAt:          resp.CreatedAt,
		UpdatedAt:          resp.UpdatedAt,
	}
}

func convertAnswerResult(resp *quizdomain.AnswerResult) quizdto.AnswerResponse {
	if resp == nil {
		return quizdto.AnswerResponse{}
	}

	return quizdto.AnswerResponse{
		QuizID:              resp.QuizID,
		SelectedAnswerIndex: resp.SelectedAnswerIndex,
		IsCorrect:           resp.IsCorrect,
		CorrectAnswerIndex:  resp.CorrectAnswerIndex,
		CorrectOptionID:     resp.CorrectOptionID,
		Explanation:         convertLocale(resp.Explanation),
	}
}

func convertLocaleRequest(src quizdto.Locale) quizdomain.Locale {
	return quizdomain.Locale{
		EN: src.EN,
		RU: src.RU,
		KK: src.KK,
	}
}

func convertLocale(src quizdomain.Locale) quizdto.Locale {
	return quizdto.Locale{
		EN: src.EN,
		RU: src.RU,
		KK: src.KK,
	}
}
