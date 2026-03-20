package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		log.Printf(
			"request_id=%s method=%s path=%s status=%d duration=%s remote=%s",
			c.GetHeader("X-Request-ID"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(startedAt).String(),
			c.ClientIP(),
		)
	}
}
