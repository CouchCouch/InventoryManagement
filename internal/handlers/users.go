package handlers

import (
	"inventory/internal/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *APIHandler) GetUserHandler(c *gin.Context) {
	id := c.Query("id")
	var items *[]domain.Item
	var err error
	if id != "" {
		ids := strings.Split(id, ",")
		items, err = s.db.ItemsByIDs(ids)
	} else {
		items, err = s.db.Items()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
