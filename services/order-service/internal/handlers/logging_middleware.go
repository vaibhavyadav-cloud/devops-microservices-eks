package handlers

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Structured JSON request logging — same intent as the Auth (Node.js) and
// Notification (Spring Boot) services' structured logging, so all services'
// logs land in the K8s log pipeline in a consistent, parseable shape.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		slog.Info("http_request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}
}
