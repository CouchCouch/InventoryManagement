package handlers

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"inventory/internal/db"
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
		items, err = s.getItemsWithBuilder(typeParam, nameParam, sortParam, limitParam, offsetParam)
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

// getItemsWithBuilder builds a query using the SafeQueryBuilder for advanced filtering/sorting
func (s *APIHandler) getItemsWithBuilder(typeParam, nameParam, sortParam, limitParam, offsetParam string) (*[]domain.Item, error) {
	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := db.NewSafeQueryBuilder(db.ItemsRegistry, selectCols)
	builder.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	builder.SetBaseWhere("i.deleted = false")

	// Apply filters
	if typeParam != "" {
		var err error
		builder, err = builder.Filter("t.name", db.OpEqual, typeParam)
		if err != nil {
			return nil, err
		}
	}

	if nameParam != "" {
		var err error
		builder, err = builder.Filter("i.name", db.OpLike, "%"+nameParam+"%")
		if err != nil {
			return nil, err
		}
	}

	// Apply sorting (format: "field:asc" or "field:desc", default: date_purchased:desc)
	sortField := "i.date_purchased"
	sortDir := db.Desc
	if sortParam != "" {
		parts := strings.Split(sortParam, ":")
		if len(parts) == 2 {
			field := parts[0]
			direction := strings.ToLower(parts[1])

			// Map user-friendly field names to actual column names
			fieldMap := map[string]string{
				"id":             "i.id",
				"name":           "i.name",
				"type":           "t.name",
				"date_purchased": "i.date_purchased",
			}

			if mapped, exists := fieldMap[field]; exists {
				sortField = mapped
			}

			if direction == "asc" {
				sortDir = db.Asc
			} else if direction == "desc" {
				sortDir = db.Desc
			}
		}
	}

	var err error
	builder, err = builder.Sort(sortField, sortDir)
	if err != nil {
		return nil, err
	}

	// Apply pagination with safe defaults
	limit := 100
	if limitParam != "" {
		if l, parseErr := strconv.Atoi(limitParam); parseErr == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}
	builder.Limit(limit)

	if offsetParam != "" {
		if offset, parseErr := strconv.Atoi(offsetParam); parseErr == nil && offset >= 0 {
			builder.Offset(offset)
		}
	}

	query, params := builder.Build()
	slog.Debug("Items query with builder", "query", query, "type", typeParam, "name", nameParam, "sort", sortParam, "limit", limit)

	rows, err := s.db.DB.Query(query, params...)
	if err != nil {
		slog.Error("Query execution failed", "error", err, "query", query)
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Item, 0)
	for rows.Next() {
		var id, name, itemType, notes string
		var datePurchased sql.NullTime
		err := rows.Scan(&id, &name, &itemType, &notes, &datePurchased)
		if err != nil {
			return nil, err
		}

		item := domain.Item{
			ID:    id,
			Name:  name,
			Type:  itemType,
			Notes: notes,
		}

		if datePurchased.Valid {
			item.DatePurchased = datePurchased.Time.Format("02-01-2006")
		}

		items = append(items, item)
	}

	return &items, nil
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
