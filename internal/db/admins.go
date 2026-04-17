package db

import (
	"database/sql"
	"encoding/base64"
	"errors"

	"inventory/internal/domain"

	"golang.org/x/crypto/argon2"
)

const (
	disableAdminByRoleQuery = `UPDATE admins SET active = FALSE WHERE role = $1;`

	createAdminQuery = `
	INSERT INTO admins (user_id, role, password_hash, active) VALUES ($1, $2, $3, TRUE)
	ON CONFLICT (user_id) DO UPDATE SET role = $2, password_hash = $3, active = TRUE;
	`

	deactivateAdminQuery = `UPDATE admins SET active = FALSE WHERE user_id = $1;`

	addAdminSessionQuery = `INSERT INTO session (session_id, user_id, expires_at) VALUES ($1, $2, $3)`

	getAdminQuery = `
	SELECT
		a.user_id,
		u.name,
		u.email,
		a.role
	FROM admins a
	JOIN users u ON a.user_id = u.id
	WHERE u.email = $1;
	`

	getPassword = `
	SELECT
		a.password_hash
	FROM admins a
	JOIN users u ON a.user_id = u.id
	WHERE u.email = $1;
	`
)

// MakeUserAdmin promotes a regular user to an admin with the specified role and password or changes the user to be another role
func (d *DB) MakeUserAdmin(admin domain.Admin) error {
	userID, err := d.CreateUser(&admin.User)
	if err != nil {
		return err
	}

	// CHANGE ME: Use a proper salt and parameters
	hashedPassword := argon2.IDKey([]byte(admin.Password), []byte{'r', 'f', 'c'}, 1, 64*1024, 4, 32)

	stringHash := base64.RawStdEncoding.EncodeToString(hashedPassword)
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec(disableAdminByRoleQuery, admin.Role); err != nil {
		tx.Rollback()
		return err
	}
	if _, err = tx.Exec(createAdminQuery, userID.String(), admin.Role, stringHash); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (d *DB) DeactivateAdmin(admin domain.Admin) error {
	_, err := d.DB.Exec(deactivateAdminQuery, admin.User.ID)
	return err
}

func (d *DB) Login(admin domain.AdminLoginRequest) error {
	hashedPassword := argon2.IDKey([]byte(admin.Password), []byte{'r', 'f', 'c'}, 1, 64*1024, 4, 32)
	var actualHash []byte

	row := d.DB.QueryRow(getPassword, admin.Email)
	err := row.Scan(&actualHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}

	stringHash := base64.RawStdEncoding.EncodeToString(hashedPassword)
	sentinel := true
	for i := range stringHash {
		if stringHash[i] != actualHash[i] {
			sentinel = false
		}
	}
	if sentinel {
		return nil
	} else {
		return domain.ErrWrongPassword
	}
}

func (d *DB ) AdminByEmail(email string) (*domain.Admin, error) {
	row := d.DB.QueryRow(getAdminQuery, email)

	var admin domain.Admin
	if err := row.Scan(&admin.User.ID, &admin.User.Name, &admin.User.Email, &admin.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &admin, nil
}
