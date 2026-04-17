package review

import (
	"curriculum-service/internal/domain"
	"curriculum-service/internal/domain/review"
	dtoreview "curriculum-service/internal/http/dto/review"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

func (h *Handler) GetReviewByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid review id")
		return
	}
	result, err := h.client.GetReviewByID(c.Request.Context(), id)
	if err != nil {
		writeReviewError(c, err)
		return
	}
	respond.JSON(c, http.StatusOK, convertReview(result))
}

func (h *Handler) CreateReview(c *gin.Context) {
	request := dtoreview.ReviewRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}
	err := h.validate.Struct(&request)
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.client.CreateReview(c.Request.Context(), convertReviewRequest(request))
	if err != nil {
		writeReviewError(c, err)
		return
	}
	respond.JSON(c, http.StatusOK, convertReview(result))
}
func (h *Handler) UpdateReview(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid review id")
		return
	}
	request := dtoreview.UpdateReviewRequest{}
	if err = c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}

	if err = h.validate.Struct(request); err != nil {
		respond.JSON(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.client.UpdateReview(c.Request.Context(), id, convertUpdateReviewRequest(request))
	if err != nil {
		writeReviewError(c, err)
		return
	}
	respond.JSON(c, http.StatusOK, convertReview(result))
}

func (h *Handler) DeleteReview(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid review id")
		return
	}
	err = h.client.DeleteReview(c.Request.Context(), id)
	if err != nil {
		writeReviewError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) GetAllReviewsByCourseID(c *gin.Context) {
	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid course id")
		return
	}
	result, err := h.client.GetAllReviewsByCourseID(c.Request.Context(), courseID)
	if err != nil {
		log.Printf("GetAllReviewsByCourseID error: %v", err)
		respond.JSON(c, http.StatusInternalServerError, "failed to get reviews")
		return
	}
	respond.JSON(c, http.StatusOK, convertReviews(result))
}

func convertReviewRequest(resp dtoreview.ReviewRequest) *review.CourseReview {
	return &review.CourseReview{
		CourseID: resp.CourseID,
		UserID:   resp.UserID,
		Rating:   resp.Rating,
		Comment:  resp.Comment,
	}
}
func convertUpdateReviewRequest(resp dtoreview.UpdateReviewRequest) *review.CourseReview {
	return &review.CourseReview{
		Rating:  resp.Rating,
		Comment: resp.Comment,
	}
}
func convertReview(resp *review.CourseReview) dtoreview.ReviewResponse {
	if resp == nil {
		return dtoreview.ReviewResponse{}
	}
	return dtoreview.ReviewResponse{
		ID:        resp.ID,
		CourseID:  resp.CourseID,
		UserID:    resp.UserID,
		Rating:    resp.Rating,
		Comment:   resp.Comment,
		ViewCount: resp.ViewCount,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}

}

func convertReviews(resp []review.CourseReview) []dtoreview.ReviewResponse {
	response := make([]dtoreview.ReviewResponse, 0, len(resp))
	for i := 0; i < len(resp); i++ {
		response = append(response, convertReview(&resp[i]))
	}
	return response
}

func writeReviewError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request query")
	default:
		c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}
