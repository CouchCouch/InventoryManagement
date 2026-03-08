package db

import (
	"database/sql"
	"errors"
	"time"

	"inventory/internal/domain"

	"github.com/google/uuid"
)

const (
	getCheckoutsQuery = `
	SELECT
		c.id,
		c.checkout_date,
		c.notes,
		u.name,
		u.email,
		c.created_by,
		a.name,
		a.email
	FROM checkouts c
	JOIN users u ON c.user_id = u.id
	JOIN users a ON c.user_id = a.id
	ORDER BY c.checkout_date DESC;
	`

	getCheckoutByIDQuery = `
	SELECT
		c.id,
		c.checkout_date,
		c.notes,
		u.id as user_id,
		u.last_name,
		u.email,
		c.created_by,
		i.id as item_id,
		i.name as item_name,
		it.name as item_type,
		ci.return_date
	FROM checkouts c
	LEFT JOIN checkout_items ci ON c.id = ci.checkout_id
	LEFT JOIN items i ON ci.item_id = i.id
	LEFT JOIN item_types it ON i.item_type_id = it.id
	WHERE c.id = $1;
	`

	createCheckoutQuery  = `INSERT INTO checkouts (user_id, notes, created_by) VALUES ($1, $2, $3) RETURNING id;`
	addCheckoutItemQuery = `INSERT INTO checkout_items (checkout_id, item_id) VALUES ($1, $2);`

	updateCheckoutQuery = `UPDATE checkouts SET notes = $1 WHERE id = $2;`

	returnItemQuery = `UPDATE checkout_items SET return_date = CURRENT_TIMESTAMP WHERE checkout_id = $1 AND item_id = $2 AND return_date IS NULL;`

	getItemStatusQuery = `
	SELECT returned FROM checkout_items
	ci JOIN eheckouts c ON c.id = ci.checkout_id
	WHERE item_id = $1 ORDER BY c.checkout_date DESC LIMIT 1
	`
)

func (d *DB) Checkouts() (*[]domain.Checkout, error) {
	rows, err := d.DB.Query(getCheckoutsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkouts := []domain.Checkout{}
	for rows.Next() {
		var checkoutID int
		var checkoutDate time.Time
		var checkoutNotes sql.NullString
		var userName, userEmail string
		var createdByName, createdByEmail string

		err := rows.Scan(
			&checkoutID,
			&checkoutDate,
			&checkoutNotes,
			&userName,
			&userEmail,
			&createdByName,
			&createdByEmail,
		)
		if err != nil {
			return nil, err
		}
		checkouts = append(checkouts, domain.Checkout{
			ID: checkoutID,
			User: domain.User{
				Name: userName,
				Email: userEmail,
			},
			CreatedBy: domain.User{
				Name: createdByName,
				Email: createdByEmail,
			},
			Notes: checkoutNotes.String,
			CheckoutDate: checkoutDate,
		})
	}

	return &checkouts, nil
}

func (d *DB) Checkout(id int) (*domain.Checkout, error) {
	rows, err := d.DB.Query(getCheckoutByIDQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkout *domain.Checkout
	for rows.Next() {
		var checkoutID int
		var checkoutDate time.Time
		var checkoutNotes sql.NullString
		var userID uuid.UUID
		var userName, userEmail string
		var createdBy uuid.UUID
		var itemID, itemName, itemType sql.NullString
		var returnDate sql.NullTime

		err := rows.Scan(
			&checkoutID,
			&checkoutDate,
			&checkoutNotes,
			&userID,
			&userName,
			&userEmail,
			&createdBy,
			&itemID,
			&itemName,
			&itemType,
			&returnDate,
		)
		if err != nil {
			return nil, err
		}

		if checkout == nil {
			checkout = &domain.Checkout{
				ID:           checkoutID,
				CheckoutDate: checkoutDate,
				Notes:        checkoutNotes.String,
				User: domain.User{
					ID:    userID,
					Name:  userName,
					Email: userEmail,
				},
				Items: []domain.CheckoutItem{},
			}
		}

		if itemID.Valid {
			checkoutItem := domain.CheckoutItem{
				Item: domain.Item{
					ID:   itemID.String,
					Name: itemName.String,
					Type: itemType.String,
				},
			}
			if returnDate.Valid {
				checkoutItem.ReturnDate = returnDate.Time
			}
			checkout.Items = append(checkout.Items, checkoutItem)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if checkout == nil {
		return nil, sql.ErrNoRows
	}

	return checkout, nil
}

func (d *DB) CreateCheckout(checkout *domain.Checkout) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}

	var checkoutID int
	err = tx.QueryRow(createCheckoutQuery, checkout.User.ID, checkout.Notes, checkout.CreatedBy).Scan(&checkoutID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, checkoutItem := range checkout.Items {
		if _, err := tx.Exec(addCheckoutItemQuery, checkoutID, checkoutItem.Item.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *DB) UpdateCheckout(checkout *domain.Checkout) error {
	_, err := d.DB.Exec(updateCheckoutQuery, checkout.Notes, checkout.ID)
	return err
}

func (d *DB) ReturnItem(checkoutID int, itemID string) error {
	res, err := d.DB.Exec(returnItemQuery, checkoutID, itemID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // Or a custom error like "item already returned or not found"
	}
	return nil
}


func (d *DB) ItemsStatus(ids []string) (*[]domain.ItemStatusResponse, error) {
	tx, err := d.DB.Begin()
	if err != nil {
		return nil, err
	}

	invalid_id := false
	statuses := make([]domain.ItemStatusResponse, 0, len(ids))

	for _, id := range(ids) {
		var status bool
		row := tx.QueryRow(getItemStatusQuery, id)
		err = row.Scan(&status)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				invalid_id = true
			} else {
				tx.Rollback()
				return nil, err
			}
		} else {
			statuses = append(statuses, domain.ItemStatusResponse{
				ID: id,
				CheckedOut: status,
			})
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	if invalid_id {
		return &statuses, domain.ErrInvalidItemID
	}
	return &statuses, nil
}
