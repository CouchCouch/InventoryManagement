package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"inventory/internal/domain"

	"github.com/gin-gonic/gin"
)

func (s *APIHandler) GetUserHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	email := c.Query("email")
	if strings.TrimSpace(email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email query parameter is required"})
		return
	}

	slog.Debug("Fetching user by email", "email", email)

	var user *domain.User
	var err error
	user, err = s.db.UserByEmail(ctx, email)
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
	c.JSON(http.StatusOK, &domain.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (s *APIHandler) CreateUserHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var user domain.User
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		slog.Warn("Invalid request body for user creation", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	slog.Info("Creating new user", "email", user.Email, "name", user.Name)

	userID, err := s.db.CreateUser(ctx, &user)
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
	c.JSON(http.StatusOK, "user created successfully")
}

func (s *APIHandler) GetMeHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
	defer cancel()

	email := c.GetString("user_email")

	var admin *domain.Admin
	var err error
	admin, err = s.db.AdminByEmail(ctx, email)
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
	c.JSON(http.StatusOK, &domain.AdminResponse{
		User: domain.UserResponse{
			ID:    admin.User.ID,
			Name:  admin.User.Name,
			Email: admin.User.Email,
		},
		Role: admin.Role,
	})
}
