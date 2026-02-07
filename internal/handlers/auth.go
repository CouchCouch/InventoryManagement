package handlers

import (
	"errors"
	"inventory/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)


func (s *APIHandler) LoginHandler(c *gin.Context) {
	var admin domain.Admin
	err := c.ShouldBindBodyWithJSON(&admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Malformed Body"})
		return
	}
	err = s.db.Login(admin)
	if err != nil {
		if errors.Is(err, domain.ErrWrongPassword) {
			c.JSON(http.StatusOK, gin.H{"login": "Wrong Password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"login": "success"})
}

func (s *APIHandler) LogoutHandler(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Login Not Implemented"})
}
