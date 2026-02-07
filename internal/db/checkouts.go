package db

import (
	"database/sql"
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
		u.id as user_id,
		u.name,
		u.email,
		c.created_by
	FROM checkouts c
	JOIN users u ON c.user_id = u.id
	ORDER BY c.id DESC;
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
)

func (d *DB) Checkouts() (*[]domain.Checkout, error) {
	rows, err := d.DB.Query(getCheckoutsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkoutsMap := make(map[int]*domain.Checkout)
	for rows.Next() {
		var checkoutID int
		var checkoutDate time.Time
		var checkoutNotes sql.NullString
		var userID uuid.UUID
		var userName, userEmail string
		var createdBy int
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

		checkout, exists := checkoutsMap[checkoutID]
		if !exists {
			checkout = &domain.Checkout{
				ID:           checkoutID,
				CheckoutDate: checkoutDate,
				Notes:        checkoutNotes.String,
				CreatedBy:    createdBy,
				User: domain.User{
					ID:        userID,
					Name:  userName,
					Email:     userEmail,
				},
				Items: []domain.CheckoutItem{},
			}
			checkoutsMap[checkoutID] = checkout
		}
	}

	checkouts := make([]domain.Checkout, 0, len(checkoutsMap))
	for _, checkout := range checkoutsMap {
		checkouts = append(checkouts, *checkout)
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
		var createdBy int
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
				CreatedBy:    createdBy,
				User: domain.User{
					ID:        userID,
					Name:  userName,
					Email:     userEmail,
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
