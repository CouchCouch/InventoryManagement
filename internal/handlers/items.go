package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"inventory/internal/domain"

	"github.com/gin-gonic/gin"
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
		if errors.Is(err, domain.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (s *APIHandler) AddItemHandler(c *gin.Context) {
	item := domain.Item{}
	err := c.ShouldBindJSON(&item)
	if err != nil {
		slog.Error("Failed to deserialize json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if item.DatePurchased != "" {
		if _, err := time.Parse("02-01-2006", item.DatePurchased); err != nil {
			slog.Error("Failed to parse time", "error", err, "date", item.DatePurchased)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	err = s.db.AddItem(&item)
	if err != nil {
		if errors.Is(err, domain.ErrItemAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		slog.Error("Failed to add item", "error", err, "item_id", item.ID)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (s *APIHandler) DeleteItemHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	} else if strings.Contains(id, ",") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulk Deletes Not supported"})
		return
	}
	err := s.db.DeleteItem(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (s *APIHandler) GetItemsStatusHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ids"})
		return
	}
	ids := strings.Split(id, ",")
	statuses, err := s.db.ItemsStatus(ids)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidItemID) {
			c.JSON(http.StatusMultiStatus, statuses)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, statuses)
}
