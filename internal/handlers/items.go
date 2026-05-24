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
	// Query parameters
	idParam := c.Query("id")
	typeParam := c.Query("type")
	nameParam := c.Query("name")
	sortParam := c.Query("sort") // Format: "field:asc" or "field:desc"
	limitParam := c.Query("limit")
	offsetParam := c.Query("offset")

	var items *[]domain.Item
	var err error

	// If using legacy id parameter, use optimized path
	if idParam != "" {
		ids := strings.Split(idParam, ",")
		items, err = s.db.ItemsByIDs(ids)
	} else if typeParam != "" && nameParam == "" && sortParam == "" && limitParam == "" && offsetParam == "" {
		// Simple type filter without sorting/pagination - use optimized path
		items, err = s.db.ItemsByType(typeParam)
	} else if typeParam == "" && nameParam == "" && sortParam == "" && limitParam == "" && offsetParam == "" {
		// No filters - use simple query
		items, err = s.db.Items()
	} else {
		// Use query builder for advanced filtering/sorting/pagination
		items, err = s.db.GetItemsWithBuilder(typeParam, nameParam, sortParam, limitParam, offsetParam)
	}

	if err != nil {
		if errors.Is(err, domain.ErrItemNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Error("Failed to query items", "error", err)
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
	c.JSON(http.StatusOK, gin.H{"id": item.ID})
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

func (s *APIHandler) GetItemsTypes(c *gin.Context) {
	types, err := s.db.ItemTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
	c.JSON(http.StatusOK, types)
}
