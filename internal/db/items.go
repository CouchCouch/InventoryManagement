package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"inventory/internal/domain"

	"github.com/lib/pq"
)

const (
	// id, name, notes, item_type, date_purchased
	addItemQuery = `INSERT into items (id, name, notes, item_type_id, date_purchased) VALUES ($1, $2, $3, $4, $5)`
	// name, notes, item_type, id
	updateItemQuery = `UPDATE items SET name = $1, notes = $2, item_type_id = $3 WHERE id = $4`
	deleteItemQuery = `UPDATE items SET deleted = true WHERE id = $1`

	selectItemTypesQuery = `SELECT name FROM item_types;`
	getTypeIDQuery       = `SELECT id FROM item_types WHERE name LIKE $1;`
	insertItemTypeQuery  = `
	INSERT INTO item_types (name) VALUES ($1)
	ON CONFLICT (name) DO UPDATE SET name = excluded.name RETURNING id;
	`
)

var itemIDRegex = regexp.MustCompile(`^[a-zA-Z0-9]{2}-[a-zA-Z0-9]{2}-[a-zA-Z0-9]{2}$`)

func (d *DB) Items(ctx context.Context) (*[]domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := NewSafeQueryBuilder(ItemsRegistry, selectCols)
	builder.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	builder.SetBaseWhere("i.deleted = false")
	if _, err := builder.Sort("i.date_purchased", Desc); err != nil {
		return nil, domain.ErrInvalidSortField
	}

	query, params := builder.Build()

	rows, err := d.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	items := make([]domain.Item, 0)
	for rows.Next() {
		var id, name, itemType, notes string
		var datePurchased sql.NullTime
		err := rows.Scan(&id, &name, &itemType, &notes, &datePurchased)
		if err != nil {
			return nil, err
		}
		if datePurchased.Valid {
			items = append(items, domain.Item{
				ID:            id,
				Name:          name,
				Type:          itemType,
				Notes:         notes,
				DatePurchased: datePurchased.Time.Format("2006-01-02"),
			})
		} else {
			items = append(items, domain.Item{
				ID:    id,
				Name:  name,
				Type:  itemType,
				Notes: notes,
			})
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &items, nil
}

func (d *DB) ItemsByIDs(ctx context.Context, ids []string) (*[]domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if len(ids) == 0 {
		return &[]domain.Item{}, nil
	}

	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := NewSafeQueryBuilder(ItemsRegistry, selectCols).
		AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id").
		SetBaseWhere("i.deleted = false")

	// Convert string slice to interface slice for IN operator
	idParams := make([]any, len(ids))
	for i, id := range ids {
		idParams[i] = id
	}
	if _, err := builder.Filter("i.id", OpIn, idParams); err != nil {
		return nil, domain.ErrInvalidFilterField
	}
	if _, err := builder.Sort("i.date_purchased", Desc); err != nil {
		return nil, domain.ErrInvalidSortField
	}

	query, params := builder.Build()
	slog.Debug("ItemsByIDs query", "query", query, "ids_count", len(ids))

	rows, err := d.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	items := make([]domain.Item, 0, len(ids))
	for rows.Next() {
		var id, name, itemType, notes string
		var datePurchased sql.NullTime
		err := rows.Scan(&id, &name, &itemType, &notes, &datePurchased)
		if err != nil {
			return nil, err
		}
		if datePurchased.Valid {
			items = append(items, domain.Item{
				ID:            id,
				Name:          name,
				Type:          itemType,
				Notes:         notes,
				DatePurchased: datePurchased.Time.Format("2006-01-02"),
			})
		} else {
			items = append(items, domain.Item{
				ID:    id,
				Name:  name,
				Type:  itemType,
				Notes: notes,
			})
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &items, nil
}

func (d *DB) Item(ctx context.Context, id string) (*domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if !validateItemID(id) {
		return nil, errors.New(domain.ErrCodeInvalidItemID)
	}

	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := NewSafeQueryBuilder(ItemsRegistry, selectCols)
	builder.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	if _, err := builder.Filter("i.id", OpEqual, id); err != nil {
		return nil, domain.ErrInvalidFilterField
	}

	query, params := builder.Build()
	row := d.DB.QueryRowContext(ctx, query, params...)

	var itemID, name, itemType, notes string
	var datePurchased sql.NullTime
	err := row.Scan(&itemID, &name, &itemType, &notes, &datePurchased)
	if err != nil {
		return nil, err
	}
	var item *domain.Item
	if datePurchased.Valid {
		item = &domain.Item{
			ID:            itemID,
			Name:          name,
			Type:          itemType,
			Notes:         notes,
			DatePurchased: datePurchased.Time.Format("2006-01-02"),
		}
	} else {
		item = &domain.Item{
			ID:    itemID,
			Name:  name,
			Type:  itemType,
			Notes: notes,
		}
	}

	return item, nil
}

func (d *DB) ItemsByType(ctx context.Context, typeFilter string) (*[]domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := NewSafeQueryBuilder(ItemsRegistry, selectCols)
	builder.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	builder.SetBaseWhere("i.deleted = false")
	if _, err := builder.Filter("t.name", OpLike, typeFilter); err != nil {
		return nil, domain.ErrInvalidFilterField
	}
	if _, err := builder.Sort("i.date_purchased", Desc); err != nil {
		return nil, domain.ErrInvalidSortField
	}

	query, params := builder.Build()
	slog.Debug("ItemsByType query", "query", query, "typeFilter", typeFilter)

	rows, err := d.DB.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	items := make([]domain.Item, 0)
	for rows.Next() {
		var id, name, itemType, notes string
		var datePurchased sql.NullTime
		err := rows.Scan(&id, &name, &itemType, &notes, &datePurchased)
		if err != nil {
			return nil, err
		}
		if datePurchased.Valid {
			items = append(items, domain.Item{
				ID:            id,
				Name:          name,
				Type:          itemType,
				Notes:         notes,
				DatePurchased: datePurchased.Time.Format("2006-01-02"),
			})
		} else {
			items = append(items, domain.Item{
				ID:    id,
				Name:  name,
				Type:  itemType,
				Notes: notes,
			})
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &items, nil
}

func (d *DB) ItemTypes(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	rows, err := d.DB.QueryContext(ctx, selectItemTypesQuery)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var itemTypes []string
	for rows.Next() {
		var itemType string
		if err := rows.Scan(&itemType); err != nil {
			return nil, err
		}
		itemTypes = append(itemTypes, itemType)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return itemTypes, nil
}

// GetItemsWithBuilder builds a query using the SafeQueryBuilder for advanced filtering/sorting
func (d *DB) GetItemsWithBuilder(ctx context.Context, typeParam, nameParam, sortParam, limitParam, offsetParam string) (*[]domain.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := NewSafeQueryBuilder(ItemsRegistry, selectCols)
	builder.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	builder.SetBaseWhere("i.deleted = false")

	// Apply filters
	if typeParam != "" {
		if _, err := builder.Filter("t.name", OpEqual, typeParam); err != nil {
			return nil, err
		}
	}

	if nameParam != "" {
		if _, err := builder.Filter("i.name", OpLike, "%"+nameParam+"%"); err != nil {
			return nil, err
		}
	}

	// Apply sorting (format: "field:asc" or "field:desc", default: date_purchased:desc)
	sortField := "i.date_purchased"
	sortDir := Desc
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

			switch direction {
			case "asc":
				sortDir = Asc
			case "desc":
				sortDir = Desc
			default:
				sortDir = Asc
			}
		}
	}

	if _, err := builder.Sort(sortField, sortDir); err != nil {
		return nil, domain.ErrInvalidSortField
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

	rows, err := d.DB.QueryContext(ctx, query, params...)
	if err != nil {
		slog.Error("Query execution failed", "error", err, "query", query)
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

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
			item.DatePurchased = datePurchased.Time.Format("2006-01-02")
		}

		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &items, nil
}

func (d *DB) itemTypeID(ctx context.Context, itemType string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	row := d.DB.QueryRowContext(ctx, getTypeIDQuery, itemType)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, errors.New(domain.ErrCodeInvalidItemType)
		}
		return -1, err
	}
	return id, nil
}

func (d *DB) addItemType(ctx context.Context, itemType string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if itemType == "" {
		return -1, errors.New("null or empty item type is not allowed")
	}
	row := d.DB.QueryRowContext(ctx, insertItemTypeQuery, itemType)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, domain.ErrInvalidItemType
		}
		return -1, err
	}
	return id, nil
}

func (d *DB) AddItem(ctx context.Context, item *domain.Item) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	slog.Debug("Adding item", "itemId", item.ID, "name", item.Name, "type", item.Type)

	if !validateItemID(item.ID) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	if item.Name == "" || item.Type == "" {
		return errors.New("item name and type cannot be empty")
	}
	// Assuming item type ID is fetched from the database or passed in some way
	itemTypeID, err := d.addItemType(ctx, item.Type)
	if err != nil {
		slog.Error("Failed to add/get item type", "error", err, "type", item.Type)
		return err
	}
	date := sql.NullTime{}
	if item.DatePurchased != "" {
		date.Time, err = time.Parse("2006-01-02", item.DatePurchased)
		if err != nil {
			date = sql.NullTime{Valid: false}
		}
	}
	_, err = d.DB.ExecContext(ctx, addItemQuery, item.ID, item.Name, item.Notes, itemTypeID, date)
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			slog.Warn("Item already exists", "itemId", item.ID)
			return domain.ErrItemAlreadyExists
		}
		slog.Error("Failed to add item", "error", err, "itemId", item.ID)
		return err
	}

	slog.Info("Item added successfully", "itemId", item.ID, "name", item.Name)
	return nil
}

func (d *DB) UpdateItem(ctx context.Context, item *domain.Item) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if !validateItemID(item.ID) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	if item.Name == "" || item.Type == "" {
		return errors.New("item name and type cannot be empty")
	}
	itemTypeID, err := d.itemTypeID(ctx, item.Type)
	if errors.Is(err, sql.ErrNoRows) {
		itemTypeID, err = d.addItemType(ctx, item.Type)
		if err != nil {
			return err
		}
	}
	_, err = d.DB.ExecContext(ctx, updateItemQuery, item.Name, item.Notes, itemTypeID, item.ID)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) DeleteItem(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if !validateItemID(id) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	_, err := d.DB.ExecContext(ctx, deleteItemQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func validateItemID(id string) bool {
	return itemIDRegex.MatchString(id)
}
