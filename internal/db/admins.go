package db

import (
	"encoding/base64"
	"inventory/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
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

	getAdmin = `
	SELECT
		a.user_id,
		u.name,
		u.email,
		a.role,
		a.password_hash,
		a.active
	FROM admins a
	JOIN users u ON a.user_id = u.id
	WHERE a.user_id = $1;
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
		logrus.Error(err)
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

func (d *DB) Login(admin domain.Admin) error {
	hashedPassword := argon2.IDKey([]byte(admin.Password), []byte{'r', 'f', 'c'}, 1, 64*1024, 4, 32)
	var actualHash []byte

	row := d.DB.QueryRow(getPassword, admin.User.Email)
	err := row.Scan(&actualHash)
	if err != nil {
		return err
	}

	stringHash := base64.RawStdEncoding.EncodeToString(hashedPassword)
	sentinel := true
	for i, _ := range stringHash {
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

var jwtKey = []byte("rfc")
var RefreshSecret = []byte("rfc-rerfresh")

func GenerateJWT(username string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 1).Unix(),
    })
    return token.SignedString(jwtKey)
}

func GenerateRefreshToken(username string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
    })
    return token.SignedString(RefreshSecret)
}
