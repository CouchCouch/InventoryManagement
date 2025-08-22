package db

import (
	"database/sql"
	"errors"
	"regexp"

	"inventory/internal/domain"
)

const (
	getItemsQuery = `
	SELECT
		i.id,
		i.name,
		i.identifiers,
		t.name,
	FROM items as i
	LEFT JOIN item_types t ON i.item_type_id = t.id;
	`

	getItemByIDQuery = `
	SELECT
		name,
		t.name,
		identifiers
	FROM items
	JOIN item_types t ON items.item_type_id = t.id
	WHERE id = $1;
	`

	addItemQuery    = `INSERT into items (name, identifiers, item_type_id) VALUES ($1, $2, $3)`
	updateItemQuery = `UPDATE items SET name = $1, identifiers = $2, item_type_id = $3 WHERE id = $4`
	deleteItemQuery = `UPDATE items SET deleted = true WHERE id = $1`

	selectItemTypesQuery = `SELECT name FROM item_types;`
	getTypeIDQuery       = `SELECT id FROM item_types WHERE name = $1;`
	insertItemTypeQuery  = `INSERT INTO item_types (name) VALUES ($1) RETURNING id`
)

var itemIDRegex = regexp.MustCompile(`^[a-zA-Z0-9]{2}-[a-zA-Z0-9]{2}-[a-zA-Z0-9]{2}$`)

func (d *db) Items() (*[]domain.Item, error) {
	rows, err := d.DB.Query(getItemsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Item, 0)
	for rows.Next() {
		var id, name, item_type, identifiers string
		err := rows.Scan(&id, &name, &item_type, &identifiers)
		if err != nil {
			return nil, err
		}
		items = append(items, domain.Item{
			ID:          id,
			Name:        name,
			Type:        item_type,
			Identifiers: identifiers,
		})
	}
	return &items, nil
}

func (d *db) Item(id string) (*domain.Item, error) {
	if !validateItemID(id) {
		return nil, errors.New(domain.ErrCodeInvalidItemID)
	}
	row := d.DB.QueryRow(getItemByIDQuery, id)
	var name, item_type, identifiers string
	err := row.Scan(&name, &item_type, &identifiers)
	if err != nil {
		return nil, err
	}
	item := &domain.Item{
		ID:          id,
		Name:        name,
		Type:        item_type,
		Identifiers: identifiers,
	}

	return item, nil
}

func (d *db) ItemTypes() ([]string, error) {
	rows, err := d.DB.Query(selectItemTypesQuery)
	if err != nil {
		return nil, err
	}
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

func (d *db) getItemTypeID(itemType string) (int, error) {
	row := d.DB.QueryRow(getTypeIDQuery, itemType)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New(domain.ErrCodeInvalidItemType)
		}
		return 0, err
	}
	return id, nil
}

func (d *db) AddItemType(itemType string) (int, error) {
	if itemType == "" {
		return 0, errors.New(domain.ErrCodeInvalidItemType)
	}
	row := d.DB.QueryRow(insertItemTypeQuery, itemType)
	var id int
	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New(domain.ErrCodeInvalidItemType)
		}
		return 0, err
	}

	return id, nil
}

func (d *db) AddItem(item *domain.Item) error {
	if !validateItemID(item.ID) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	if item.Name == "" || item.Type == "" {
		return errors.New("item name and type cannot be empty")
	}
	// Assuming item type ID is fetched from the database or passed in some way
	itemTypeID, err := d.getItemTypeID(item.Type)
	if errors.Is(err, sql.ErrNoRows) {
		itemTypeID, err = d.AddItemType(item.Type)
		if err != nil {
			return err
		}
	}
	_, err = d.DB.Exec(addItemQuery, item.Name, item.Identifiers, itemTypeID)
	if err != nil {
		return err
	}
	return nil
}

func (d *db) UpdateItem(item *domain.Item) error {
	if !validateItemID(item.ID) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	if item.Name == "" || item.Type == "" {
		return errors.New("item name and type cannot be empty")
	}
	itemTypeID, err := d.getItemTypeID(item.Type)
	if errors.Is(err, sql.ErrNoRows) {
		itemTypeID, err = d.AddItemType(item.Type)
		if err != nil {
			return err
		}
	}
	_, err = d.DB.Exec(updateItemQuery, item.Name, item.Identifiers, itemTypeID, item.ID)
	if err != nil {
		return err
	}
	return nil
}

func (d *db) DeleteItem(id string) error {
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
