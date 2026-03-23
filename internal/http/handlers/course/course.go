package course

import (
	"curriculum-service/internal/domain"
	domaincourse "curriculum-service/internal/domain/course"
	dtocourse "curriculum-service/internal/http/dto/course"
	"curriculum-service/internal/http/dto/durationcategory"
	"curriculum-service/internal/http/dto/level"
	"curriculum-service/internal/http/dto/status"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) ListCourses(c *gin.Context) {
	result, err := h.client.GetAllCourses(c.Request.Context())
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertCourses(result))
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

func convertCourses(resp []domaincourse.Courses) []dtocourse.Courses {
	courses := make([]dtocourse.Courses, len(resp))

	for i := range resp {
		roles := make([]dtocourse.Role, len(resp[i].Author.Roles))
		for j := range resp[i].Author.Roles {
			roles[j] = dtocourse.Role{
				ID:           resp[i].Author.Roles[j].ID,
				Code:         resp[i].Author.Roles[j].Code,
				Name:         resp[i].Author.Roles[j].Name,
				Description:  resp[i].Author.Roles[j].Description,
				IsDefault:    resp[i].Author.Roles[j].IsDefault,
				IsPrivileged: resp[i].Author.Roles[j].IsPrivileged,
				IsSupport:    resp[i].Author.Roles[j].IsSupport,
				CreatedAt:    resp[i].Author.Roles[j].CreatedAt,
			}
		}

		courses[i] = dtocourse.Courses{
			ID:             resp[i].ID,
			ExpectedHours:  resp[i].ExpectedHours,
			Rating:         resp[i].Rating,
			RatingCount:    resp[i].RatingCount,
			StudentsCount:  resp[i].StudentsCount,
			LessonsCount:   resp[i].LessonsCount,
			HasCertificate: resp[i].HasCertificate,
			CoverImageUrl:  resp[i].CoverImageUrl,
			PublishedAt:    resp[i].PublishedAt,
			CreatedAt:      resp[i].CreatedAt,
			UpdatedAt:      resp[i].UpdatedAt,
			Status: status.Status{
				ID:   resp[i].Status.ID,
				Name: resp[i].Status.Name,
				Code: resp[i].Status.Code,
			},
			Level: level.Level{
				ID:   resp[i].Level.ID,
				Name: resp[i].Level.Name,
				Code: resp[i].Level.Code,
			},
			DurationCategory: durationcategory.DurationCategory{
				ID:   resp[i].DurationCategory.ID,
				Name: resp[i].DurationCategory.Name,
				Code: resp[i].DurationCategory.Code,
			},
			Author: dtocourse.User{
				ID:        resp[i].Author.ID,
				Email:     resp[i].Author.Email,
				Roles:     roles,
				IsActive:  resp[i].Author.IsActive,
				CreatedAt: resp[i].Author.CreatedAt,
			},
		}
	}

	return courses
}
