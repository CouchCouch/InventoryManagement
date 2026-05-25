package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"inventory/internal/domain"

	"github.com/gin-gonic/gin"
)

func (s *APIHandler) LoginHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var admin domain.AdminLoginRequest
	err := c.ShouldBindBodyWithJSON(&admin)
	if err != nil {
		slog.Warn("Login request with invalid JSON", "ip", c.ClientIP(), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"login": "fail"})
		return
	}

	slog.Info("Login attempt", "email", admin.Email, "ip", c.ClientIP())

	err = s.db.Login(ctx, admin)
	if err != nil {
		if errors.Is(err, domain.ErrWrongPassword) {
			slog.Warn("Login failed - wrong password", "email", admin.Email, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"login": "fail"})
			return
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			slog.Warn("Login failed - user not found", "email", admin.Email, "ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{"login": "fail"})
			return
		}
		slog.Error("Login failed - database error", "error", err, "email", admin.Email, "ip", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"login": "fail"})
		return
	}

	tokenString, err := s.auth.GenerateAccessToken(admin.Email)
	if err != nil {
		slog.Error("Error generating JWT", "error", err, "email", admin.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"login": "fail"})
		return
	}

	slog.Info("Login successful", "email", admin.Email, "ip", c.ClientIP())
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"login": "success"})
}

func (s *APIHandler) LogoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{})
}
