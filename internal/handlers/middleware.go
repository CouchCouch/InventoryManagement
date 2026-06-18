package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *APIHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			slog.Warn("Unauthorized access attempt - missing token", "ip", c.ClientIP(), "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token", "message": "Unauthorized"})
			c.Abort()
			return
		}

		email, err := s.auth.ValidateToken(tokenString)
		if err != nil {
			slog.Warn("Invalid token attempt", "ip", c.ClientIP(), "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_email", email)
		c.Next()
	}
}

func RequestLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		path := c.Request.URL.Path
		method := c.Request.Method

		slog.Info("Request started",
			"request_id", requestID,
			"method", method,
			"path", path,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent())

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logAttrs := []any{
			"request_id", requestID,
			"method", method,
			"path", path,
			"status", statusCode,
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
		}

		if email, exists := c.Get("user_email"); exists {
			logAttrs = append(logAttrs, "user_email", email)
		}

		switch {
		case statusCode >= 500:
			slog.Error("Request completed with server error", logAttrs...)
		case statusCode >= 400:
			slog.Warn("Request completed with client error", logAttrs...)
		default:
			slog.Info("Request completed", logAttrs...)
		}
	}
}
