package db

import (
	"database/sql"
	"errors"

	"inventory/internal/domain"
)

const (
	getUsersQuery = `
	SELECT
		id,
		first_name,
		last_name,
		email
	FROM users
	WHERE deleted = FALSE;
	`

	getUserByIDQuery = `
	SELECT
		id,
		first_name,
		last_name,
		email
	FROM users
	WHERE id = $1 AND deleted = FALSE;
	`

	createUserQuery = `INSERT INTO users (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING id;`
	updateUserQuery = `UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4;`
	deleteUserQuery = `UPDATE users SET deleted = TRUE WHERE id = $1;`
)

func (d *db) Users() (*[]domain.User, error) {
	rows, err := d.DB.Query(getUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}

func (d *db) User(id int) (*domain.User, error) {
	row := d.DB.QueryRow(getUserByIDQuery, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (d *db) CreateUser(user *domain.User) (int, error) {
	var id int
	err := d.DB.QueryRow(createUserQuery, user.FirstName, user.LastName, user.Email).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *db) UpdateUser(user *domain.User) error {
	_, err := d.DB.Exec(updateUserQuery, user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (d *db) DeleteUser(id int) error {
	_, err := d.DB.Exec(deleteUserQuery, id)
	if err != nil {
		return err
	}
	return nil
}
