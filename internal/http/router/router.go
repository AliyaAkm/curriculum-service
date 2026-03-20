package router

import (
	"curriculum-service/internal/http/handlers"

	"github.com/gin-gonic/gin"
)

func New(catalogH *handlers.CatalogHandler, globalMiddlewares []gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(globalMiddlewares...)

	r.GET("/health", health)

	r.GET("/courses", catalogH.ListCourses)
	r.GET("/courses/search", catalogH.SearchCourses)
	r.GET("/courses/filter", catalogH.FilterCourses)
	r.GET("/courses/filters", catalogH.ListFilterOptions)

	return r
}
