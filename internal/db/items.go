package db

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"

	"inventory/internal/domain"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	getItemsQuery = `
	SELECT
		i.id,
		i.name,
		t.name,
		i.notes,
		i.date_purchased
	FROM items AS i
	LEFT JOIN item_types t ON i.item_type_id = t.id
	ORDER BY i.date_purchased DESC;
	`

	getItemByIDQuery = `
	SELECT
		i.id,
		i.name,
		t.name,
		i.notes,
		i.date_purchased
	FROM items AS i
	LEFT JOIN item_types t ON i.item_type_id = t.id
	WHERE id = $1;
	`

	getItemsByIDsQuery = `
	SELECT
		i.id,
		i.name,
		t.name,
		i.notes,
		i.date_purchased
	FROM items AS i
	LEFT JOIN item_types t ON i.item_type_id = t.id
	WHERE ($1) IN (i.id)
	ORDER BY i.date_purchased DESC;
	`

	addItemQuery    = `INSERT into items (id, name, notes, item_type_id, date_purchased) VALUES ($1, $2, $3, $4, $5)`
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
	rows, err := d.DB.Query(getItemsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Item, 0)
	for rows.Next() {
		var id, name, item_type, notes string
		var date_purchased sql.NullTime
		err := rows.Scan(&id, &name, &item_type, &notes, &date_purchased)
		if err != nil {
			return nil, err
		}
		if date_purchased.Valid {
			items = append(items, domain.Item{
				ID:            id,
				Name:          name,
				Type:          item_type,
				Notes:         notes,
				DatePurchased: date_purchased.Time.Format("02-01-2006"),
			})
		} else {
			items = append(items, domain.Item{
				ID:            id,
				Name:          name,
				Type:          item_type,
				Notes:         notes,
			})
		}
	}
	return &items, nil
}

func (d *DB) ItemsByIDs(ids []string) (*[]domain.Item, error) {
	idq := strings.Join(ids, ",")
	rows, err := d.DB.Query(getItemsByIDsQuery, idq)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Item, 0, len(ids))
	for rows.Next() {
		var id, name, item_type, notes string
		var date_purchased sql.NullTime
		err := rows.Scan(&id, &name, &item_type, &notes, &date_purchased)
		if err != nil {
			return nil, err
		}
		log.Info("Date: ", date_purchased)
		if date_purchased.Valid {
			items = append(items, domain.Item{
				ID:            id,
				Name:          name,
				Type:          item_type,
				Notes:         notes,
				DatePurchased: date_purchased.Time.Format("02-01-2006"),
			})
		} else {
			items = append(items, domain.Item{
				ID:            id,
				Name:          name,
				Type:          item_type,
				Notes:         notes,
			})
		}
	}

	return &items, nil
}

func (d *DB) Item(id string) (*domain.Item, error) {
	if !validateItemID(id) {
		return nil, errors.New(domain.ErrCodeInvalidItemID)
	}
	row := d.DB.QueryRow(getItemByIDQuery, id)
	var name, item_type, notes string
	var date_purchased sql.NullTime
	err := row.Scan(&name, &item_type, &notes, &date_purchased)
	if err != nil {
		return nil, err
	}
	item := &domain.Item{}
	log.Info("Date: ", date_purchased)
	if date_purchased.Valid {
		item = &domain.Item{
			ID:            id,
			Name:          name,
			Type:          item_type,
			Notes:         notes,
			DatePurchased: date_purchased.Time.Format("02-01-2006"),
		}
	} else {
		item = &domain.Item{
			ID:            id,
			Name:          name,
			Type:          item_type,
			Notes:         notes,
		}
	}

	return item, nil
}

func (d *DB) ItemTypes() ([]string, error) {
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

func (d *DB) getItemTypeID(itemType string) (int, error) {
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
	if !validateItemID(item.ID) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	if item.Name == "" || item.Type == "" {
		return errors.New("item name and type cannot be empty")
	}
	// Assuming item type ID is fetched from the database or passed in some way
	itemTypeID, err := d.addItemType(item.Type)
	if err != nil {
		return err
	}
	_, err = d.DB.Exec(addItemQuery, item.ID, item.Name, item.Notes, itemTypeID, item.DatePurchased)
	if err != nil {
		if (err.(*pq.Error).Code == "23505") {
			return domain.ErrItemAlreadyExists
		}
	}
	return err
}

func (d *DB) UpdateItem(item *domain.Item) error {
	if !validateItemID(item.ID) {
		return errors.New(domain.ErrCodeInvalidItemID)
	}
	if item.Name == "" || item.Type == "" {
		return errors.New("item name and type cannot be empty")
	}
	itemTypeID, err := d.getItemTypeID(item.Type)
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
