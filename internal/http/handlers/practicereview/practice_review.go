package practicereview

import (
	"curriculum-service/internal/domain"
	practicereviewdomain "curriculum-service/internal/domain/practicereview"
	practicereviewdto "curriculum-service/internal/http/dto/practicereview"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) CreateSubmission(c *gin.Context) {
	claims, ok := h.claims(c)
	if !ok {
		return
	}
	studentID, ok := userIDFromClaims(c, claims)
	if !ok {
		return
	}

	practiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "validation", "invalid practice id")
		return
	}

	var req practicereviewdto.CreateSubmissionRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	result, err := h.client.CreateSubmission(c.Request.Context(), practicereviewdomain.CreateSubmissionRequest{
		PracticeID: practiceID,
		StudentID:  studentID,
		Code:       req.Code,
		Language:   req.Language,
		Output:     req.Output,
		Error:      req.Error,
		ErrorType:  req.ErrorType,
		ExitCode:   req.ExitCode,
		DurationMS: req.DurationMS,
	})
	if err != nil {
		writeReviewError(c, err)
		return
	}

	respond.JSON(c, http.StatusCreated, submissionResponse(result))
}

func (h *Handler) ListMySubmissions(c *gin.Context) {
	claims, ok := h.claims(c)
	if !ok {
		return
	}
	studentID, ok := userIDFromClaims(c, claims)
	if !ok {
		return
	}

	filter, ok := studentFilter(c, studentID)
	if !ok {
		return
	}

	items, err := h.client.ListStudentSubmissions(c.Request.Context(), filter)
	if err != nil {
		writeReviewError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, submissionResponses(items))
}

func (h *Handler) GetMySubmission(c *gin.Context) {
	claims, ok := h.claims(c)
	if !ok {
		return
	}
	studentID, ok := userIDFromClaims(c, claims)
	if !ok {
		return
	}

	submissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "validation", "invalid submission id")
		return
	}

	result, err := h.client.GetStudentSubmission(c.Request.Context(), studentID, submissionID)
	if err != nil {
		writeReviewError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, submissionResponse(result))
}

func (h *Handler) ListTeacherSubmissions(c *gin.Context) {
	claims, ok := h.claims(c)
	if !ok {
		return
	}
	teacherID, ok := userIDFromClaims(c, claims)
	if !ok {
		return
	}
	if !isTeacherOrAdmin(claims) {
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
		return
	}

	filter, ok := teacherFilter(c, teacherID, middleware.ClaimsHasRole(claims, middleware.RoleAdmin))
	if !ok {
		return
	}

	items, err := h.client.ListTeacherSubmissions(c.Request.Context(), filter)
	if err != nil {
		writeReviewError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, submissionResponses(items))
}

func (h *Handler) GetTeacherSubmission(c *gin.Context) {
	claims, ok := h.claims(c)
	if !ok {
		return
	}
	teacherID, ok := userIDFromClaims(c, claims)
	if !ok {
		return
	}
	if !isTeacherOrAdmin(claims) {
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
		return
	}

	submissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "validation", "invalid submission id")
		return
	}

	result, err := h.client.GetTeacherSubmission(c.Request.Context(), teacherID, middleware.ClaimsHasRole(claims, middleware.RoleAdmin), submissionID)
	if err != nil {
		writeReviewError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, submissionResponse(result))
}

func (h *Handler) ReviewSubmission(c *gin.Context) {
	claims, ok := h.claims(c)
	if !ok {
		return
	}
	teacherID, ok := userIDFromClaims(c, claims)
	if !ok {
		return
	}
	if !isTeacherOrAdmin(claims) {
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
		return
	}

	submissionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.Error(c, http.StatusBadRequest, "validation", "invalid submission id")
		return
	}

	var req practicereviewdto.ReviewSubmissionRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	result, err := h.client.ReviewSubmission(c.Request.Context(), practicereviewdomain.ReviewSubmissionRequest{
		SubmissionID: submissionID,
		TeacherID:    teacherID,
		IsAdmin:      middleware.ClaimsHasRole(claims, middleware.RoleAdmin),
		Status:       req.Status,
		Comment:      req.Comment,
	})
	if err != nil {
		writeReviewError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, submissionResponse(result))
}

func (h *Handler) claims(c *gin.Context) (*middleware.Claims, bool) {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return nil, false
	}
	return claims, true
}

func userIDFromClaims(c *gin.Context, claims *middleware.Claims) (uuid.UUID, bool) {
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respond.Error(c, http.StatusUnauthorized, "unauthorized", "invalid user")
		return uuid.Nil, false
	}
	return userID, true
}

func studentFilter(c *gin.Context, studentID uuid.UUID) (practicereviewdomain.StudentListFilter, bool) {
	filter := practicereviewdomain.StudentListFilter{
		StudentID: studentID,
		Status:    c.Query("status"),
	}

	if courseID, ok := optionalUUIDQuery(c, "course_id"); !ok {
		return practicereviewdomain.StudentListFilter{}, false
	} else {
		filter.CourseID = courseID
	}

	if practiceID, ok := optionalUUIDQuery(c, "practice_id"); !ok {
		return practicereviewdomain.StudentListFilter{}, false
	} else {
		filter.PracticeID = practiceID
	}

	return filter, true
}

func teacherFilter(c *gin.Context, teacherID uuid.UUID, isAdmin bool) (practicereviewdomain.TeacherListFilter, bool) {
	filter := practicereviewdomain.TeacherListFilter{
		TeacherID: teacherID,
		IsAdmin:   isAdmin,
		Status:    c.Query("status"),
		Limit:     intQuery(c, "limit", 50),
		Offset:    intQuery(c, "offset", 0),
	}

	if courseID, ok := optionalUUIDQuery(c, "course_id"); !ok {
		return practicereviewdomain.TeacherListFilter{}, false
	} else {
		filter.CourseID = courseID
	}

	if practiceID, ok := optionalUUIDQuery(c, "practice_id"); !ok {
		return practicereviewdomain.TeacherListFilter{}, false
	} else {
		filter.PracticeID = practiceID
	}

	if studentID, ok := optionalUUIDQuery(c, "student_id"); !ok {
		return practicereviewdomain.TeacherListFilter{}, false
	} else {
		filter.StudentID = studentID
	}

	return filter, true
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

func isTeacherOrAdmin(claims *middleware.Claims) bool {
	return middleware.ClaimsHasRole(claims, middleware.RoleTeacher) || middleware.ClaimsHasRole(claims, middleware.RoleAdmin)
}

func submissionResponses(items []practicereviewdomain.Submission) []practicereviewdto.SubmissionResponse {
	response := make([]practicereviewdto.SubmissionResponse, len(items))
	for i := range items {
		response[i] = submissionResponse(&items[i])
	}
	return response
}

func submissionResponse(src *practicereviewdomain.Submission) practicereviewdto.SubmissionResponse {
	if src == nil {
		return practicereviewdto.SubmissionResponse{}
	}

	return practicereviewdto.SubmissionResponse{
		ID:                    src.ID,
		PracticeID:            src.PracticeID,
		StudentID:             src.StudentID,
		CourseID:              src.CourseID,
		LessonID:              src.LessonID,
		Status:                src.Status,
		Code:                  src.Code,
		Language:              src.Language,
		Output:                src.Output,
		Error:                 src.Error,
		ErrorType:             src.ErrorType,
		ExitCode:              src.ExitCode,
		DurationMS:            src.DurationMS,
		TeacherComment:        src.TeacherComment,
		ReviewedBy:            src.ReviewedBy,
		ReviewedAt:            src.ReviewedAt,
		AttemptNumber:         src.AttemptNumber,
		PracticeTitle:         src.PracticeTitle,
		StudentEmail:          src.StudentEmail,
		CourseTitle:           src.CourseTitle,
		LessonTitle:           src.LessonTitle,
		ProgressStatus:        src.ProgressStatus,
		ProgressStartedAt:     src.ProgressStartedAt,
		ProgressCompletedAt:   src.ProgressCompletedAt,
		ProgressLastAttemptAt: src.ProgressLastAttemptAt,
		ProgressAttemptsCount: src.ProgressAttemptsCount,
		CreatedAt:             src.CreatedAt,
		UpdatedAt:             src.UpdatedAt,
	}
}

func writeReviewError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation), errors.Is(err, domain.ErrInvalidPracticeReviewStatus), errors.Is(err, domain.ErrPracticeManualReviewNotAllowed):
		respond.Error(c, http.StatusBadRequest, "validation", err.Error())
	case errors.Is(err, domain.ErrPracticePrerequisitesNotMet):
		respond.Error(c, http.StatusConflict, "practice_prerequisites_not_met", err.Error())
	case errors.Is(err, domain.ErrPracticeNotFound), errors.Is(err, domain.ErrPracticeSubmissionNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, domain.ErrForbidden):
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
	case errors.Is(err, domain.ErrCourseSubscriptionNotFound):
		respond.Error(c, http.StatusForbidden, "course_subscription_not_found", domain.ErrCourseSubscriptionNotFound.Error())
	case errors.Is(err, domain.ErrPracticeSubmissionExists), errors.Is(err, domain.ErrPracticeAlreadyCompleted):
		respond.Error(c, http.StatusConflict, "conflict", err.Error())
	default:
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}
