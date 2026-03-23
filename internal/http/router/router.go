package router

import (
	"curriculum-service/internal/http/handlers/course"
	"curriculum-service/internal/http/handlers/durationcategory"
	"curriculum-service/internal/http/handlers/level"
	"curriculum-service/internal/http/handlers/status"
	"curriculum-service/internal/http/handlers/tag"
	"curriculum-service/internal/http/handlers/topic"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Status           *status.Handler
	Level            *level.Handler
	DurationCategory *durationcategory.Handler
	Topic            *topic.Handler
	Tag              *tag.Handler
	Course           *course.Handler
}

func New(handler Handler, globalMiddlewares []gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(globalMiddlewares...)

	r.GET("/health", health)

	r.GET("/course", handler.Course.ListCourses)
	r.GET("/dictionary/status", handler.Status.ListCourseStatuses)
	r.GET("/dictionary/level", handler.Level.ListCourseLevels)
	r.GET("/dictionary/duration_category", handler.DurationCategory.ListCourseDurationCategories)
	r.GET("/dictionary/topic", handler.Topic.ListCourseTopics)
	r.GET("/dictionary/tag", handler.Tag.ListCourseTags)
	//r.GET("/course/search", catalogH.SearchCourses)
	//r.GET("/course/filter", catalogH.FilterCourses)
	//r.GET("/course/filters", catalogH.ListFilterOptions)

	return r
}
