package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS(origins, methods, headers []string) gin.HandlerFunc {
	allowedOrigins := normalizeValues(origins)
	allowedMethods := normalizeValues(methods)
	allowedHeaders := normalizeValues(headers)

	allowAllOrigins := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"
	methodsValue := strings.Join(allowedMethods, ", ")
	headersValue := strings.Join(allowedHeaders, ", ")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		switch {
		case allowAllOrigins:
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		case origin != "" && contains(allowedOrigins, origin):
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", methodsValue)
		c.Writer.Header().Set("Access-Control-Allow-Headers", headersValue)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func normalizeValues(values []string) []string {
	items := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
