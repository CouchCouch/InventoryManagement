package handlers

import (
	"errors"
	"net/http"

	"inventory/internal/domain"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *APIHandler) LoginHandler(c *gin.Context) {
	var admin domain.Admin
	err := c.ShouldBindBodyWithJSON(&admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	err = s.db.Login(admin)
	if err != nil {
		if errors.Is(err, domain.ErrWrongPassword) {
			c.JSON(http.StatusOK, gin.H{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	tokenString, err := s.auth.GenerateJWT(admin.User.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		log.Error("Error generating JWT: ", err)
		return
	}
	c.Header("Authorization", "Bearer "+tokenString)
	c.JSON(http.StatusOK, gin.H{})
}
