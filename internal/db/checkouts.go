package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"inventory/internal/domain"
)

const (
	createCheckoutQuery  = `INSERT INTO checkouts (user_id, notes, created_by, checkout_date) VALUES ($1, $2, $3, $4) RETURNING id;`
	addCheckoutItemQuery = `INSERT INTO checkout_items (checkout_id, item_id) VALUES ($1, $2);`

	updateCheckoutQuery = `UPDATE checkouts SET notes = $1 WHERE id = $2;`

	returnItemQuery = `UPDATE checkout_items SET return_date = CURRENT_TIMESTAMP WHERE checkout_id = $1 AND item_id = $2 AND return_date IS NULL;`

	getItemStatusQuery = `
	SELECT EXISTS (
	    SELECT 1 FROM checkout_items ci
	    WHERE ci.item_id = $1 AND ci.return_date IS NULL
	)
	`
)

func (d *DB) Checkouts(ctx context.Context) ([]domain.Checkout, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	selectCols := `c.id, c.checkout_date, c.notes, u.name, u.email, a.name, a.email`
	builder := NewSafeQueryBuilder(CheckoutsRegistry, selectCols)
	builder.AddJoin("JOIN users u ON c.user_id = u.id")
	builder.AddJoin("JOIN users a ON c.created_by = a.id")
	_, err := builder.Sort("c.checkout_date", Asc)
	if err != nil {
		return nil, domain.ErrInvalidSortField
	}

	query, params := builder.Build()

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

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
				Name:  userName,
				Email: userEmail,
			},
			CreatedBy: domain.User{
				Name:  createdByName,
				Email: createdByEmail,
			},
			Notes:        checkoutNotes.String,
			CheckoutDate: checkoutDate,
		})
	}

	ids := make([]int, len(checkouts))
	for i, c := range checkouts {
		ids[i] = c.ID
	}

	itemsByID, err := checkoutItemsByCheckoutIDs(ctx, tx, ids)
	if err != nil {
		slog.Debug("Error getting checkout items", "err", err)
	}
	for i := range checkouts {
		if items, ok := itemsByID[checkouts[i].ID]; ok {
			checkouts[i].Items = items
		} else {
			checkouts[i].Items = []domain.CheckoutItem{}
		}
	}

	return checkouts, nil
}

func checkoutItemsByCheckoutIDs(ctx context.Context, tx *sql.Tx, ids []int) (map[int][]domain.CheckoutItem, error) {
	if len(ids) == 0 {
		return map[int][]domain.CheckoutItem{}, nil
	}

	selectCols := `ci.checkout_id, ci.item_id, ci.return_date, i.name, t.name, i.notes`
	builder := NewSafeQueryBuilder(CheckoutItemsRegistry, selectCols)
	builder.AddJoin("JOIN items i ON ci.item_id = i.id")
	builder.AddJoin("JOIN item_types t ON i.item_type_id = t.id")

	idParams := make([]any, len(ids))
	for i, id := range ids {
		idParams[i] = id
	}
	if _, err := builder.Filter("ci.checkout_id", OpIn, idParams); err != nil {
		return nil, err
	}
	if _, err := builder.Sort("ci.item_id", Asc); err != nil {
		return nil, domain.ErrInvalidSortField
	}

	query, params := builder.Build()

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	itemsByID := make(map[int][]domain.CheckoutItem)
	for rows.Next() {
		var checkoutID int
		var itemID, itemName, itemType string
		var notes sql.NullString
		var returnDate sql.NullTime

		err := rows.Scan(
			&checkoutID,
			&itemID,
			&returnDate,
			&itemName,
			&itemType,
			&notes,
		)
		if err != nil {
			return nil, err
		}

		item := domain.CheckoutItem{
			Item: domain.Item{
				ID:   itemID,
				Name: itemName,
				Type: itemType,
			},
			ReturnDate: returnDate.Time,
		}
		if notes.Valid {
			item.Item.Notes = notes.String
		}

		itemsByID[checkoutID] = append(itemsByID[checkoutID], item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return itemsByID, nil
}

func checkoutItems(ctx context.Context, tx *sql.Tx, id int) ([]domain.CheckoutItem, error) {
	selectCols := `ci.item_id, ci.return_date, i.name, t.name, i.notes`
	builder := NewSafeQueryBuilder(CheckoutItemsRegistry, selectCols)
	builder.AddJoin("JOIN items i ON ci.item_id = i.id")
	builder.AddJoin("JOIN item_types t ON i.item_type_id = t.id")
	if _, err := builder.Filter("ci.checkout_id", OpEqual, id); err != nil {
		return nil, domain.ErrInvalidFilterField
	}
	_, err := builder.Sort("ci.item_id", Asc)
	if err != nil {
		return nil, domain.ErrInvalidSortField
	}
	query, params := builder.Build()

	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	checkoutItems := []domain.CheckoutItem{}
	for rows.Next() {
		var itemID, itemName, itemType string
		var notes sql.NullString
		var returnDate sql.NullTime

		err := rows.Scan(
			&itemID,
			&returnDate,
			&itemName,
			&itemType,
			&notes,
		)
		if err != nil {
			return nil, err
		}

		if notes.Valid {
			checkoutItems = append(checkoutItems, domain.CheckoutItem{
				Item: domain.Item{
					ID:    itemID,
					Name:  itemName,
					Type:  itemType,
					Notes: notes.String,
				},
				ReturnDate: returnDate.Time,
			})
		} else {
			checkoutItems = append(checkoutItems, domain.CheckoutItem{
				Item: domain.Item{
					ID:   itemID,
					Name: itemName,
					Type: itemType,
				},
				ReturnDate: returnDate.Time,
			})
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return checkoutItems, nil
}

func (d *DB) Checkout(ctx context.Context, id int) (*domain.Checkout, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	selectCols := `c.id, c.checkout_date, c.notes, u.name, u.email, a.name, a.email`
	builder := NewSafeQueryBuilder(CheckoutsRegistry, selectCols)
	builder.AddJoin("JOIN users u ON c.user_id = u.id")
	builder.AddJoin("JOIN users a ON c.created_by = a.id")
	_, err := builder.Sort("c.checkout_date", Asc)
	if err != nil {
		return nil, domain.ErrInvalidSortField
	}

	if _, err = builder.Filter("c.id", OpEqual, id); err != nil {
		return nil, err
	}

	query, params := builder.Build()

	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, query, params...)

	var checkoutID int
	var checkoutDate time.Time
	var checkoutNotes sql.NullString
	var userName, userEmail string
	var createdByName, createdByEmail string

	err = row.Scan(
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

	checkout := domain.Checkout{
		ID: checkoutID,
		User: domain.User{
			Name:  userName,
			Email: userEmail,
		},
		CreatedBy: domain.User{
			Name:  createdByName,
			Email: createdByEmail,
		},
		Notes:        checkoutNotes.String,
		CheckoutDate: checkoutDate,
	}

	checkoutItems, err := checkoutItems(ctx, tx, checkoutID)
	if err != nil {
		checkout.Items = []domain.CheckoutItem{}
	} else {
		checkout.Items = checkoutItems
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &checkout, nil
}

func (d *DB) CreateCheckout(ctx context.Context, user domain.User, items []string, checkoutDate time.Time, createdBy domain.Admin, notes string) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		slog.Error("Failed to begin transaction", "error", err)
		return 0, err
	}

	//nolint:errcheck
	defer tx.Rollback()

	var checkoutID int
	err = tx.QueryRowContext(ctx, createCheckoutQuery, user.ID, notes, createdBy.User.ID, checkoutDate).Scan(&checkoutID)
	if err != nil {
		slog.Error("Failed to create checkout", "error", err)
		return 0, err
	}

	for _, id := range items {
		if _, err := tx.ExecContext(ctx, addCheckoutItemQuery, checkoutID, id); err != nil {
			slog.Error("Failed to add checkout item", "error", err, "checkout_id", checkoutID, "item_id", id)
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return checkoutID, nil
}

func (d *DB) UpdateCheckout(ctx context.Context, checkout *domain.Checkout) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	_, err := d.DB.ExecContext(ctx, updateCheckoutQuery, checkout.Notes, checkout.ID)
	return err
}

func (d *DB) ReturnItem(ctx context.Context, checkoutID int, items []string) error {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer tx.Rollback()

	for _, itemID := range items {
		if _, err := tx.ExecContext(ctx, returnItemQuery, checkoutID, itemID); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *DB) ItemsStatus(ctx context.Context, ids []string) (*[]domain.ItemStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	//nolint
	defer tx.Rollback()

	invalidID := false
	statuses := make([]domain.ItemStatusResponse, 0, len(ids))

	for _, id := range ids {
		var status bool
		row := tx.QueryRowContext(ctx, getItemStatusQuery, id)
		err = row.Scan(&status)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				invalidID = true
			} else {
				return nil, err
			}
		} else {
			statuses = append(statuses, domain.ItemStatusResponse{
				ID:         id,
				CheckedOut: status,
			})
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	if invalidID {
		return &statuses, domain.ErrInvalidItemID
	}
	return &statuses, nil
}
