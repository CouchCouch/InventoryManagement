package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"inventory/internal/domain"

	"github.com/gin-gonic/gin"
)

func (s *APIHandler) GetUserHandler(c *gin.Context) {
	email := c.Query("email")
	if strings.TrimSpace(email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email query parameter is required"})
		return
	}
	
	slog.Debug("Fetching user by email", "email", email)
	
	var user *domain.User
	var err error
	user, err = s.db.UserByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			slog.Debug("User not found", "email", email)
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Error("Error fetching user", "error", err, "email", email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *APIHandler) CreateUserHandler(c *gin.Context) {
	var user domain.User
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		slog.Warn("Invalid request body for user creation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	
	slog.Info("Creating new user", "email", user.Email, "name", user.Name)
	
	userID, err := s.db.CreateUser(&user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			slog.Warn("User creation failed - already exists", "email", user.Email)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Error("User creation failed", "error", err, "email", user.Email)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	
	slog.Info("User created successfully", "user_id", userID, "email", user.Email)
	c.JSON(http.StatusOK, "item added successfully")
}
