package handlers

import (
	"net/http"
	"strings"

	"inventory/internal/domain"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *APIHandler) GetItemsHandler(c *gin.Context) {
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

func (s *APIHandler) AddItemHandler(c *gin.Context) {
	item := domain.Item{}
	err := c.ShouldBindJSON(&item)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = s.db.AddItem(&item)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "item added successfully")
}
