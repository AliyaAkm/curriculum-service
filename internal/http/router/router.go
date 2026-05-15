package router

import (
	"curriculum-service/internal/http/handlers/course"
	"curriculum-service/internal/http/handlers/coursepoint"
	"curriculum-service/internal/http/handlers/durationcategory"
	"curriculum-service/internal/http/handlers/lesson"
	"curriculum-service/internal/http/handlers/level"
	"curriculum-service/internal/http/handlers/locale"
	"curriculum-service/internal/http/handlers/module"
	"curriculum-service/internal/http/handlers/practice"
	"curriculum-service/internal/http/handlers/quiz"
	"curriculum-service/internal/http/handlers/review"
	"curriculum-service/internal/http/handlers/status"
	"curriculum-service/internal/http/handlers/streak"
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
	Locale           *locale.Handler
	Module           *module.Handler
	Lesson           *lesson.Handler
	Practice         *practice.Handler
	Quiz             *quiz.Handler
	Review           *review.Handler
	CoursePoint      *coursepoint.Handler
	Streak           *streak.Handler
}

func New(handler Handler, globalMiddlewares []gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(globalMiddlewares...)

	r.GET("/health", health)

	r.GET("/course", handler.Course.ListCourses)
	// cправочники
	r.GET("/dictionary/status", handler.Status.ListCourseStatuses)
	r.GET("/dictionary/level", handler.Level.ListCourseLevels)
	r.GET("/dictionary/duration_category", handler.DurationCategory.ListCourseDurationCategories)
	r.GET("/dictionary/topic", handler.Topic.ListCourseTopics)
	r.GET("/dictionary/tag", handler.Tag.ListCourseTags)
	r.GET("/dictionary/locale", handler.Locale.ListCourseLocales)

	// курсы
	r.POST("/course", handler.Course.CreateCourse)
	r.GET("/course/:id", handler.Course.GetCourseByID)
	r.DELETE("/course/:id", handler.Course.DeleteCourse)
	r.PUT("/course/:id", handler.Course.UpdateCourse)

	// course subscription
	r.POST("/course/enrollment", handler.Course.CreateSubscription)

	// модули
	r.GET("/module", handler.Module.GetAllModules)
	r.POST("/module", handler.Module.CreateModule)
	r.GET("/module/:id", handler.Module.GetModuleByID)
	r.PUT("/module/:id", handler.Module.UpdateModule)
	r.DELETE("/module/:id", handler.Module.DeleteModule)

	// lesson
	r.GET("/module/lesson/:id", handler.Lesson.GetAllLessons)
	r.GET("/lesson/:id", handler.Lesson.GetLessonByID)
	r.PUT("/lesson/:id", handler.Lesson.UpdateLesson)
	r.POST("/lesson", handler.Lesson.CreateLesson)

	// practice
	r.POST("/practice", handler.Practice.CreatePractice)
	r.GET("/practice", handler.Practice.GetPracticeByLessonID)
	r.GET("/practice/:id", handler.Practice.GetPracticeByID)
	r.PUT("/practice/:id", handler.Practice.UpdatePractice)
	r.DELETE("/practice/:id", handler.Practice.DeletePractice)

	// quiz
	r.POST("/quiz", handler.Quiz.CreateQuiz)
	r.GET("/lesson/:id/quiz", handler.Quiz.GetQuizzesByLessonID)
	r.POST("/quiz/:id/answer", handler.Quiz.SubmitAnswer)
	r.GET("/quiz/:id", handler.Quiz.GetQuizByID)
	r.PUT("/quiz/:id", handler.Quiz.UpdateQuiz)
	r.DELETE("/quiz/:id", handler.Quiz.DeleteQuiz)
	//r.GET("/course/search", catalogH.SearchCourses)
	//r.GET("/course/filter", catalogH.FilterCourses)
	//r.GET("/course/filters", catalogH.ListFilterOptions)

	// review
	r.GET("/review/:id", handler.Review.GetReviewByID)
	r.POST("/review", handler.Review.CreateReview)
	r.PUT("/review/:id", handler.Review.UpdateReview)
	r.DELETE("/review/:id", handler.Review.DeleteReview)
	r.GET("/course/review/:id", handler.Review.GetAllReviewsByCourseID)

	// coursePoint
	r.POST("/point", handler.CoursePoint.CreateCoursePoint)
	r.PUT("/point/:id", handler.CoursePoint.UpdateCoursePoint)
	r.DELETE("/point/:id", handler.CoursePoint.DeleteCoursePoint)
	r.GET("/leaderboard/:id", handler.CoursePoint.GetCoursePointByCourseID)

	// daily streak
	r.GET("/streak", handler.Streak.GetStreak)

	return r
}
