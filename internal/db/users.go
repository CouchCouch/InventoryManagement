package db

import (
	"database/sql"
	"errors"

	"inventory/internal/domain"

	"github.com/google/uuid"
)

const (
	getUsersQuery = `
	SELECT
		id,
		name,
		email
	FROM users;
	`

	getUserByIDQuery = `
	SELECT
		id,
		name,
		email
	FROM users
	WHERE id = $1;
	`

	getUserByEmailQuery = `
	SELECT
		id,
		name,
		email
	FROM users
	WHERE email = $1;
	`

	createUserQuery = `
	INSERT INTO users (id, name, email) VALUES ($1, $2, $3)
	ON CONFLICT (email) DO UPDATE SET email = excluded.email RETURNING id;
	`
)

// Users retrieves all users from the database
func (d *DB) Users() (*[]domain.User, error) {
	rows, err := d.DB.Query(getUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}

// User retrieves a user by their ID
func (d *DB) User(id int) (*domain.User, error) {
	row := d.DB.QueryRow(getUserByIDQuery, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// User retrieves a user by their Email
func (d *DB) UserByEmail(email string) (*domain.User, error) {
	row := d.DB.QueryRow(getUserByEmailQuery, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database, is there is already a user with the given email, the user id is returned
func (d *DB) CreateUser(user *domain.User) (uuid.UUID, error) {
	id := uuid.New()
	err := d.DB.QueryRow(createUserQuery, id, user.Name, user.Email).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
