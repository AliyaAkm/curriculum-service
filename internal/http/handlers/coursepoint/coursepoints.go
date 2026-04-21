package coursepoint

import (
	"curriculum-service/internal/domain"
	domaincoursepoint "curriculum-service/internal/domain/coursepoint"
	"curriculum-service/internal/http/dto/coursepoint"
	"curriculum-service/internal/http/respond"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (h *Handler) CreateCoursePoint(c *gin.Context) {
	request := coursepoint.CreateUserCoursePointsRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.client.CreateCoursePoint(c.Request.Context(), convertCoursePointRequest(request))
	if err != nil {
		writeCoursePointError(c, err)
		return
	}
	respond.JSON(c, http.StatusOK, convertCoursePoint(result))
}
func (h *Handler) UpdateCoursePoint(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid lesson id")
		return
	}
	request := coursepoint.UpdateUserCoursePointsRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.client.UpdateCoursePoint(c.Request.Context(), id, convertUpdateCoursePointRequest(request))
	if err != nil {
		writeCoursePointError(c, err)
		return
	}
	respond.JSON(c, http.StatusOK, convertCoursePoint(result))
}

func (h *Handler) DeleteCoursePoint(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid id")
		return
	}
	err = h.client.DeleteCoursePoint(c.Request.Context(), id)
	if err != nil {
		writeCoursePointError(c, err)
	}
	c.Status(http.StatusNoContent)
}
func (h *Handler) GetCoursePointByCourseID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respond.JSON(c, http.StatusBadRequest, "invalid course id")
		return
	}
	result, err := h.client.GetCoursePointByCourseID(c.Request.Context(), id)
	if err != nil {
		writeCoursePointError(c, err)
		return
	}
	respond.JSON(c, http.StatusOK, convertLeaderboard(result))
}

func convertLeaderboard(resp []domaincoursepoint.Leaderboard) []coursepoint.Leaderboard {
	leaderboard := make([]coursepoint.Leaderboard, len(resp))
	for i := range resp {
		leaderboard[i] = coursepoint.Leaderboard{
			Place:  resp[i].Place,
			UserID: resp[i].UserID,
			XP:     resp[i].XP,
		}
	}
	return leaderboard
}

func convertUpdateCoursePointRequest(resp coursepoint.UpdateUserCoursePointsRequest) *domaincoursepoint.UserCoursePoints {
	return &domaincoursepoint.UserCoursePoints{
		XP: resp.XP,
	}
}

func convertCoursePointRequest(resp coursepoint.CreateUserCoursePointsRequest) *domaincoursepoint.UserCoursePoints {
	return &domaincoursepoint.UserCoursePoints{
		LessonID: resp.LessonID,
		UserID:   resp.UserID,
		XP:       resp.XP,
	}
}
func convertCoursePoint(resp *domaincoursepoint.UserCoursePoints) coursepoint.UserCoursePointsResponse {
	return coursepoint.UserCoursePointsResponse{
		ID:        resp.ID,
		LessonID:  resp.LessonID,
		UserID:    resp.UserID,
		XP:        resp.XP,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}
}
func writeCoursePointError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		respond.Error(c, http.StatusBadRequest, "validation", "invalid request query")
	default:
		c.Error(err)
		respond.Error(c, http.StatusInternalServerError, "internal", domain.ErrInternal.Error())
	}
}
