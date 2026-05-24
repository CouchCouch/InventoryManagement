package db

import (
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

func (d *DB) Items() (*[]domain.Item, error) {
	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := NewSafeQueryBuilder(ItemsRegistry, selectCols)
	builder.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	builder.SetBaseWhere("i.deleted = false")
	if _, err := builder.Sort("i.date_purchased", Desc); err != nil {
		return nil, domain.ErrInvalidSortField
	}

	query, params := builder.Build()
	slog.Debug("Items query", "query", query, "params", params)

	rows, err := d.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer rows.Close()

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
				DatePurchased: datePurchased.Time.Format("02-01-2006"),
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
	return &items, nil
}

func (d *DB) ItemsByIDs(ids []string) (*[]domain.Item, error) {
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
	builder.builder.Filter("i.id", OpIn, idParams)
	if _, err := builder.Sort("i.date_purchased", Desc); err != nil {
		return nil, domain.ErrInvalidSortField
	}

	query, params := builder.Build()
	slog.Debug("ItemsByIDs query", "query", query, "ids_count", len(ids))

	rows, err := d.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer rows.Close()

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
				DatePurchased: datePurchased.Time.Format("02-01-2006"),
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

	return &items, nil
}

func (d *DB) Item(id string) (*domain.Item, error) {
	if !validateItemID(id) {
		return nil, errors.New(domain.ErrCodeInvalidItemID)
	}

	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := NewSafeQueryBuilder(ItemsRegistry, selectCols)
	builder.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	builder.builder.Filter("i.id", OpEqual, id)

	query, params := builder.Build()
	row := d.DB.QueryRow(query, params...)

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
			DatePurchased: datePurchased.Time.Format("02-01-2006"),
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

func (d *DB) ItemsByType(typeFilter string) (*[]domain.Item, error) {
	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := NewSafeQueryBuilder(ItemsRegistry, selectCols)
	builder.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	builder.SetBaseWhere("i.deleted = false")
	builder.builder.Filter("t.name", OpLike, "%"+typeFilter+"%")
	if _, err := builder.Sort("i.date_purchased", Desc); err != nil {
		return nil, domain.ErrInvalidSortField
	}

	query, params := builder.Build()
	slog.Debug("ItemsByType query", "query", query, "typeFilter", typeFilter)

	rows, err := d.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer rows.Close()

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
				DatePurchased: datePurchased.Time.Format("02-01-2006"),
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
	return &items, nil
}

func (d *DB) ItemTypes() ([]string, error) {
	rows, err := d.DB.Query(selectItemTypesQuery)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer rows.Close()

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
func (d *DB) GetItemsWithBuilder(typeParam, nameParam, sortParam, limitParam, offsetParam string) (*[]domain.Item, error) {
	selectCols := `i.id, i.name, t.name, i.notes, i.date_purchased`
	builder := NewSafeQueryBuilder(ItemsRegistry, selectCols)
	builder.AddJoin("LEFT JOIN item_types t ON i.item_type_id = t.id")
	builder.SetBaseWhere("i.deleted = false")

	// Apply filters
	if typeParam != "" {
		var err error
		builder, err = builder.Filter("t.name", OpEqual, typeParam)
		if err != nil {
			return nil, err
		}
	}

	if nameParam != "" {
		var err error
		builder, err = builder.Filter("i.name", OpLike, "%"+nameParam+"%")
		if err != nil {
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

	var err error
	builder, err = builder.Sort(sortField, sortDir)
	if err != nil {
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

	rows, err := d.DB.Query(query, params...)
	if err != nil {
		slog.Error("Query execution failed", "error", err, "query", query)
		return nil, err
	}

	//nolint:errcheck
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

func (d *DB) itemTypeID(itemType string) (int, error) {
	row := d.DB.QueryRow(getTypeIDQuery, itemType)
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

func (d *DB) addItemType(itemType string) (int, error) {
	if itemType == "" {
		return -1, errors.New("null or empty item type is not allowed")
	}
	row := d.DB.QueryRow(insertItemTypeQuery, itemType)
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

func (d *DB) AddItem(item *domain.Item) error {
	slog.Debug("Adding item", "itemId", item.ID, "name", item.Name, "type", item.Type)

	if !validateItemID(item.ID) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	if item.Name == "" || item.Type == "" {
		return errors.New("item name and type cannot be empty")
	}
	// Assuming item type ID is fetched from the database or passed in some way
	itemTypeID, err := d.addItemType(item.Type)
	if err != nil {
		slog.Error("Failed to add/get item type", "error", err, "type", item.Type)
		return err
	}
	date := sql.NullTime{}
	if item.DatePurchased != "" {
		date.Time, err = time.Parse("02-01-2006", item.DatePurchased)
		if err != nil {
			date = sql.NullTime{Valid: false}
		}
	}
	_, err = d.DB.Exec(addItemQuery, item.ID, item.Name, item.Notes, itemTypeID, date)
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

func (d *DB) UpdateItem(item *domain.Item) error {
	if !validateItemID(item.ID) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	if item.Name == "" || item.Type == "" {
		return errors.New("item name and type cannot be empty")
	}
	itemTypeID, err := d.itemTypeID(item.Type)
	if errors.Is(err, sql.ErrNoRows) {
		itemTypeID, err = d.addItemType(item.Type)
		if err != nil {
			return err
		}
	}
	_, err = d.DB.Exec(updateItemQuery, item.Name, item.Notes, itemTypeID, item.ID)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) DeleteItem(id string) error {
	if !validateItemID(id) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	_, err := d.DB.Exec(deleteItemQuery, id)
	if err != nil {
		return err
	}
	return nil
}

func validateItemID(id string) bool {
	return itemIDRegex.MatchString(id)
}
