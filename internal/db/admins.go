package db

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"inventory/internal/domain"
)

const (
	getAdminsQuery = `
	SELECT
		u.id,
		u.first_name,
		u.last_name,
		u.email,
		a.role
	FROM admins a
	JOIN users u ON a.user_id = u.id
	WHERE u.deleted = FALSE;
	`

	getAdminByIDQuery = `
	SELECT
		u.id,
		u.first_name,
		u.last_name,
		u.email,
		a.role
	FROM admins a
	JOIN users u ON a.user_id = u.id
	WHERE u.id = $1 AND u.deleted = FALSE;
	`

	createAdminQuery = `INSERT INTO admins (user_id, role, password_hash) VALUES ($1, $2, $3);`
	updateAdminQuery = `UPDATE admins SET role = $1, password_hash = $2 WHERE user_id = $3;`
	deleteAdminQuery = `DELETE FROM admins WHERE user_id = $1;`
)

func (d *db) Admins() (*[]domain.Admin, error) {
	rows, err := d.DB.Query(getAdminsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []domain.Admin
	for rows.Next() {
		var admin domain.Admin
		if err := rows.Scan(&admin.User.ID, &admin.User.FirstName, &admin.User.LastName, &admin.User.Email, &admin.Role); err != nil {
			return nil, err
		}
		admins = append(admins, admin)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &admins, nil
}

func (d *db) Admin(id int) (*domain.Admin, error) {
	row := d.DB.QueryRow(getAdminByIDQuery, id)

	var admin domain.Admin
	err := row.Scan(&admin.User.ID, &admin.User.FirstName, &admin.User.LastName, &admin.User.Email, &admin.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("admin not found")
		}
		return nil, err
	}
	return &admin, nil
}

func (d *db) CreateAdmin(admin *domain.Admin) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}

	userID, err := d.CreateUser(&admin.User)
	if err != nil {
		tx.Rollback()
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(createAdminQuery, userID, admin.Role, string(hashedPassword))
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (d *db) UpdateAdmin(admin *domain.Admin) error {
	tx, err := d.DB.Begin()
	if err != nil {
		return err
	}

	err = d.UpdateUser(&admin.User)
	if err != nil {
		tx.Rollback()
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(updateAdminQuery, admin.Role, string(hashedPassword), admin.User.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (d *db) DeleteAdmin(id int) error {
	_, err := d.DB.Exec(deleteAdminQuery, id)
	if err != nil {
		return err
	}
	return nil
}
