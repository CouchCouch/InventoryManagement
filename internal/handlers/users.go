package handlers

import (
	"errors"
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
	var user *domain.User
	var err error
	user, err = s.db.UserByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *APIHandler) CreateUserHandler(c *gin.Context) {
	var user domain.User
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	_, err = s.db.CreateUser(&user)
	if err != nil {
		if errors.Is(err, domain.ErrItemAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, "item added successfully")
}
