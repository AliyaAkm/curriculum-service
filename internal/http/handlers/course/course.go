package course

import (
	"curriculum-service/internal/domain"
	domaincourse "curriculum-service/internal/domain/course"
	domaintag "curriculum-service/internal/domain/tag"
	dtocourse "curriculum-service/internal/http/dto/course"
	"curriculum-service/internal/http/dto/durationcategory"
	"curriculum-service/internal/http/dto/level"
	"curriculum-service/internal/http/dto/status"
	dtotag "curriculum-service/internal/http/dto/tag"
	"curriculum-service/internal/http/dto/topic"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) ListCourses(c *gin.Context) {
	var query dtocourse.GetCoursesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid query params")
		return
	}
	result, err := h.client.GetAllCourses(c.Request.Context(), query)
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertCourses(result))
}

func (h *Handler) CreateCourse(c *gin.Context) {
	request := dtocourse.CourseRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.client.CreateCourse(c.Request.Context(), convertCourseRequest(request))
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertCourse(result))
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	request := dtocourse.SubscriptionRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.client.CreateSubscription(c.Request.Context(), convertSubscriptionRequest(request))
	if err != nil {
		writeCatalogError(c, err)
		return
	}
	respond.JSON(c, http.StatusOK, convertSubscription(result))
}

func (h *Handler) GetCourseByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid course id")
		return
	}
	result, err := h.client.GetCourseByID(c.Request.Context(), id)
	if err != nil {
		writeCatalogError(c, err)
		return
	}
	respond.JSON(c, http.StatusOK, convertCourse(result))

}

func (h *Handler) DeleteCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid course id")
		return
	}
	err = h.client.DeleteCourse(c.Request.Context(), id)
	if err != nil {
		writeCatalogError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) UpdateCourse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid course id")
		return
	}
	request := dtocourse.CourseRequest{}
	if err = c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.client.UpdateCourse(c.Request.Context(), id, convertCourseRequest(request))
	if err != nil {
		writeCatalogError(c, err)
		return
	}

	respond.JSON(c, http.StatusOK, convertCourse(result))
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

func convertCourses(resp []domaincourse.Course) []dtocourse.Courses {
	courses := make([]dtocourse.Courses, len(resp))

	for i := range resp {
		tags := make([]dtotag.Tag, len(resp[i].Tags))
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
		for j := range resp[i].Tags {
			tags[j] = dtotag.Tag{
				ID:   resp[i].Tags[j].ID,
				Code: resp[i].Tags[j].Code,
				Name: resp[i].Tags[j].Name,
			}
		}

		courses[i] = dtocourse.Courses{
			ID:             resp[i].ID,
			Title:          resp[i].Title,
			SubTitle:       resp[i].SubTitle,
			Description:    resp[i].Description,
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
			Tags: tags,
			Topic: topic.Topic{
				ID:   resp[i].Topic.ID,
				Name: resp[i].Topic.Name,
				Code: resp[i].Topic.Code,
			},
			LearningOutcomes: resp[i].LearningOutcomes,
		}
	}

	return courses
}

func convertSubscription(resp *domaincourse.Subscription) dtocourse.Subscription {
	return dtocourse.Subscription{
		ID:       resp.ID,
		UserID:   resp.UserID,
		CourseID: resp.CourseID,
	}
}

func convertCourse(resp *domaincourse.Course) dtocourse.Courses {
	tags := make([]dtotag.Tag, len(resp.Tags))
	roles := make([]dtocourse.Role, len(resp.Author.Roles))
	for j := range resp.Author.Roles {
		roles[j] = dtocourse.Role{
			ID:           resp.Author.Roles[j].ID,
			Code:         resp.Author.Roles[j].Code,
			Name:         resp.Author.Roles[j].Name,
			Description:  resp.Author.Roles[j].Description,
			IsDefault:    resp.Author.Roles[j].IsDefault,
			IsPrivileged: resp.Author.Roles[j].IsPrivileged,
			IsSupport:    resp.Author.Roles[j].IsSupport,
			CreatedAt:    resp.Author.Roles[j].CreatedAt,
		}
	}
	for j := range resp.Tags {
		tags[j] = dtotag.Tag{
			ID:   resp.Tags[j].ID,
			Code: resp.Tags[j].Code,
			Name: resp.Tags[j].Name,
		}
	}

	return dtocourse.Courses{
		ID:             resp.ID,
		Title:          resp.Title,
		SubTitle:       resp.SubTitle,
		Description:    resp.Description,
		ExpectedHours:  resp.ExpectedHours,
		Rating:         resp.Rating,
		RatingCount:    resp.RatingCount,
		StudentsCount:  resp.StudentsCount,
		LessonsCount:   resp.LessonsCount,
		HasCertificate: resp.HasCertificate,
		CoverImageUrl:  resp.CoverImageUrl,
		PublishedAt:    resp.PublishedAt,
		CreatedAt:      resp.CreatedAt,
		UpdatedAt:      resp.UpdatedAt,
		Status: status.Status{
			ID:   resp.Status.ID,
			Name: resp.Status.Name,
			Code: resp.Status.Code,
		},
		Level: level.Level{
			ID:   resp.Level.ID,
			Name: resp.Level.Name,
			Code: resp.Level.Code,
		},
		DurationCategory: durationcategory.DurationCategory{
			ID:   resp.DurationCategory.ID,
			Name: resp.DurationCategory.Name,
			Code: resp.DurationCategory.Code,
		},
		Author: dtocourse.User{
			ID:        resp.Author.ID,
			Email:     resp.Author.Email,
			Roles:     roles,
			IsActive:  resp.Author.IsActive,
			CreatedAt: resp.Author.CreatedAt,
		},
		Tags: tags,
		Topic: topic.Topic{
			ID:   resp.Topic.ID,
			Name: resp.Topic.Name,
			Code: resp.Topic.Code,
		},
		LearningOutcomes: resp.LearningOutcomes,
	}
}

func convertSubscriptionRequest(resp dtocourse.SubscriptionRequest) *domaincourse.Subscription {
	return &domaincourse.Subscription{
		UserID:   resp.UserID,
		CourseID: resp.CourseID,
	}
}

func convertCourseRequest(resp dtocourse.CourseRequest) *domaincourse.Course {
	tags := make([]domaintag.Tag, len(resp.TagIDs))
	for i, id := range resp.TagIDs {
		tags[i] = domaintag.Tag{
			ID: id,
		}
	}

	return &domaincourse.Course{
		Title:         resp.Title,
		SubTitle:      resp.SubTitle,
		Description:   resp.Description,
		ExpectedHours: resp.ExpectedHours,
		//Rating:             resp.Rating,
		//RatingCount:        resp.RatingCount,
		StudentsCount:      resp.StudentsCount,
		LessonsCount:       resp.LessonsCount,
		HasCertificate:     resp.HasCertificate,
		CoverImageUrl:      resp.CoverImageUrl,
		StatusID:           resp.StatusID,
		LevelID:            resp.LevelID,
		DurationCategoryID: resp.DurationCategoryID,
		AuthorID:           resp.AuthorID,
		TopicID:            resp.TopicID,
		Tags:               tags,
		LearningOutcomes:   resp.LearningOutcomes,
	}
}
