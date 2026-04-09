package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"inventory/internal/domain"

	"github.com/gin-gonic/gin"
)

func (s *APIHandler) LoginHandler(c *gin.Context) {
	var admin domain.Admin
	err := c.ShouldBindBodyWithJSON(&admin)
	if err != nil {
		slog.Warn("Login request with invalid JSON", "ip", c.ClientIP(), "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	
	slog.Info("Login attempt", "email", admin.User.Email, "ip", c.ClientIP())
	
	err = s.db.Login(admin)
	if err != nil {
		if errors.Is(err, domain.ErrWrongPassword) {
			slog.Warn("Login failed - wrong password", "email", admin.User.Email, "ip", c.ClientIP())
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		slog.Error("Login failed - database error", "error", err, "email", admin.User.Email, "ip", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	
	tokenString, err := s.auth.GenerateJWT(admin.User.Email)
	if err != nil {
		slog.Error("Error generating JWT", "error", err, "email", admin.User.Email)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	
	slog.Info("Login successful", "email", admin.User.Email, "ip", c.ClientIP())
	c.Header("Authorization", "Bearer "+tokenString)
	c.JSON(http.StatusOK, gin.H{})
}
