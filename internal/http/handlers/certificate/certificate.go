package certificate

import (
	"curriculum-service/internal/domain"
	certificatedomain "curriculum-service/internal/domain/certificate"
	certificatedto "curriculum-service/internal/http/dto/certificate"
	"curriculum-service/internal/http/middleware"
	"curriculum-service/internal/http/respond"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ссылка на сертификат доступен 30 минут
func (h *Handler) IssueCourseCertificate(c *gin.Context) {
	claims := middleware.GetClaims(h.jwtMgr, c)
	if claims == nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respond.JSON(c, http.StatusUnauthorized, "invalid user id")
		return
	}

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid course id")
		return
	}

	result, err := h.client.IssueCourseCertificate(
		c.Request.Context(),
		userID,
		courseID,
		claims.Login,
		middleware.ClaimsHasRole(claims, middleware.RoleAdmin),
	)
	if err != nil {
		writeCertificateError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertIssueCertificateResponse(result))
}

func writeCertificateError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCourseNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", domain.ErrCourseNotFound.Error())
	case errors.Is(err, domain.ErrCertificateNotFound):
		respond.Error(c, http.StatusNotFound, "not_found", domain.ErrCertificateNotFound.Error())
	case errors.Is(err, domain.ErrCourseNotCompleted):
		respond.Error(c, http.StatusConflict, "course_not_completed", domain.ErrCourseNotCompleted.Error())
	case errors.Is(err, domain.ErrCertificateUnavailable):
		respond.Error(c, http.StatusForbidden, "certificate_unavailable", domain.ErrCertificateUnavailable.Error())
	case errors.Is(err, domain.ErrForbidden):
		respond.Error(c, http.StatusForbidden, "forbidden", domain.ErrForbidden.Error())
	default:
		_ = c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}

func convertIssueCertificateResponse(src *certificatedomain.IssueResult) certificatedto.IssueCertificateResponse {
	if src == nil {
		return certificatedto.IssueCertificateResponse{}
	}

	return certificatedto.IssueCertificateResponse{
		CertificateNumber: src.Certificate.CertificateNumber,
		IssuedAt:          src.Certificate.IssuedAt,
		CompletedAt:       src.Certificate.CompletedAt,
		DownloadURL:       src.DownloadURL,
	}
}
