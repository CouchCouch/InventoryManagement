package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"inventory/internal/domain"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

	createUserQuery = `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`
)

// Users retrieves all users from the database
func (d *DB) Users(ctx context.Context) (*[]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := d.DB.QueryContext(ctx, getUsersQuery)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
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
func (d *DB) User(ctx context.Context, id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	row := d.DB.QueryRowContext(ctx, getUserByIDQuery, id)

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

func (d *DB) UserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	row := d.DB.QueryRowContext(ctx, getUserByEmailQuery, email)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in the database, if there is already a user with the given email, the user id is returned
func (d *DB) CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	id := uuid.New()
	_, err := d.DB.ExecContext(ctx, createUserQuery, id, user.Name, user.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				slog.Warn("User already exists", "email", user.Email)
				return uuid.Nil, domain.ErrUserAlreadyExists
			}
		}
		return uuid.Nil, err
	}

	return id, nil
}
