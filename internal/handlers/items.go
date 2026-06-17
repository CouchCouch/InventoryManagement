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

func (s *APIHandler) GetItemsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	// Query parameters
	idParam := c.Query("id")
	typeParam := c.Query("type")
	nameParam := c.Query("name")
	sortParam := c.Query("sort") // Format: "field:asc" or "field:desc"
	limitParam := c.Query("limit")
	offsetParam := c.Query("offset")

	var items *[]domain.Item
	var err error

	if idParam != "" {
		ids := strings.Split(idParam, ",")
		items, err = s.db.ItemsByIDs(ctx, ids)
	} else if typeParam != "" && nameParam == "" && sortParam == "" && limitParam == "" && offsetParam == "" {
		items, err = s.db.ItemsByType(ctx, typeParam)
	} else if typeParam == "" && nameParam == "" && sortParam == "" && limitParam == "" && offsetParam == "" {
		items, err = s.db.Items(ctx)
	} else {
		items, err = s.db.GetItemsWithBuilder(ctx, typeParam, nameParam, sortParam, limitParam, offsetParam)
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
	defer cancel()

	item := domain.Item{}
	err := c.ShouldBindJSON(&item)
	if err != nil {
		slog.Error("Failed to deserialize json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = s.db.AddItem(ctx, &item)
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
	defer cancel()

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	} else if strings.Contains(id, ",") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulk Deletes Not supported"})
		return
	}
	err := s.db.DeleteItem(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (s *APIHandler) GetItemsStatusHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ids"})
		return
	}
	ids := strings.Split(id, ",")
	statuses, err := s.db.ItemsStatus(ctx, ids)
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
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second)
	defer cancel()

	types, err := s.db.ItemTypes(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
	c.JSON(http.StatusOK, types)
}
