package db

import (
	"inventory/internal/domain"

	"golang.org/x/crypto/argon2"
)

const (
	createAdminQuery = `
	UPDATE admins SET active = FALSE where role = $2;
	INSERT INTO admins (user_id, role, password_hash, active) VALUES ($1, $2, $3, TRUE)
	ON CONFLICT (user_id) DO UPDATE SET role = $1, password_hash = $3, active = TRUE;
	`

	deactivateAdminQuery = `UPDATE admins SET active = FALSE WHERE user_id = $1;`

	addAdminSessionQuery = `INSERT INTO session (session_id, user_id, expires_at) VALUES ($1, $2, $3)`

	getAdmin = `
	SELECT
		a.user_id,
		u.first_name,
		u.last_name,
		u.email,
		a.role,
		a.password_hash,
		a.active
	FROM admins a
	JOIN users u ON a.user_id = u.id
	WHERE a.user_id = $1;
	`
)

// MakeUserAdmin promotes a regular user to an admin with the specified role and password or changes the user to be another role
func (d *DB) MakeUserAdmin(admin domain.Admin, password string) error {
	userID, err := d.CreateUser(&admin.User)
	if err != nil {
		return err
	}

	// CHANGE ME: Use a proper salt and parameters
	hashedPassword := argon2.IDKey([]byte(password), []byte{'r', 'f', 'c'}, 1, 64*1024, 4, 32)

	_, err = d.DB.Exec(createAdminQuery, userID, admin.Role, string(hashedPassword))
	return err
}


func (d *DB) DeactivateAdmin(admin domain.Admin) error {
	_, err := d.DB.Exec(deactivateAdminQuery, admin.User.ID)
	return err
}
