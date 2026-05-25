package db

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"inventory/internal/domain"

	"golang.org/x/crypto/argon2"
)

type parsedHash struct {
	Memory      uint32
	Time        uint32
	Parallelism uint8
	Salt        []byte
	Hash        []byte
}

const (
	argon2idPHCStringFormat = "$argon2id$%s9$m=%d,t=%d,p=%d$%s$%s"
	argon2idVersion = "v=19"

	disableAdminByRoleQuery = `UPDATE admins SET active = FALSE WHERE role = $1;`

	createAdminQuery = `
	INSERT INTO admins (user_id, role, password_hash, active) VALUES ($1, $2, $3, TRUE)
	ON CONFLICT (user_id) DO UPDATE SET role = $2, password_hash = $3, active = TRUE;
	`

	deactivateAdminQuery = `UPDATE admins SET active = FALSE WHERE user_id = $1;`

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

func generatePHC(password string, time, memory uint32, threads uint8, saltLen, keyLen uint32) (string, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil { return "", err }
	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	phc := fmt.Sprintf(argon2idPHCStringFormat, argon2idVersion, memory, time, threads, b64Salt, b64Hash)
	return phc, nil
}

func parsePHC(phc string) (*parsedHash, error) {
	parts := strings.Split(phc, "$")

	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != argon2idVersion {
		return nil, domain.ErrInvalidHash
	}

	params := strings.Split(parts[3], ",")
	var memory, time, parallel uint64
	var err error
	for _, p := range params {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			return nil, domain.ErrInvalidHash
		}
		switch kv[0] {
		case "m":
			memory, err = strconv.ParseUint(kv[1], 10, 32)
		case "t":
			time, err = strconv.ParseUint(kv[1], 10, 32)
		case "p":
			parallel, err = strconv.ParseUint(kv[1], 10, 8)
		default:
			err = domain.ErrInvalidHash
		}
		if err != nil {
			return nil, err
		}
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return nil, domain.ErrInvalidHash
	}
	hash, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil {
		return nil, domain.ErrInvalidHash
	}

	return &parsedHash{
		Memory:      uint32(memory),
		Time:        uint32(time),
		Parallelism: uint8(parallel),
		Salt:        salt,
		Hash:        hash,
	}, nil
}

// MakeUserAdmin promotes a regular user to an admin with the specified role and password or changes the user to be another role
func (d *DB) MakeUserAdmin(ctx context.Context, admin domain.Admin) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userID, err := d.CreateUser(ctx, &admin.User)
	if err != nil {
		return err
	}

	hashedPassword, err := generatePHC(admin.Password, 1, 64*1024, 4, 32, 32)
	if err != nil {
		return err
	}
	tx, err := d.DB.BeginTx(ctx, &sql.TxOptions{})

	//nolint:errcheck
	defer tx.Rollback()

	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, disableAdminByRoleQuery, admin.Role); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, createAdminQuery, userID.String(), admin.Role, hashedPassword); err != nil {
		return err
	}
	return tx.Commit()
}

func (d *DB) DeactivateAdmin(ctx context.Context, admin domain.Admin) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := d.DB.ExecContext(ctx, deactivateAdminQuery, admin.User.ID)
	return err
}

func (d *DB) Login(ctx context.Context, admin domain.AdminLoginRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var storedPHC string
	row := d.DB.QueryRowContext(ctx, getPassword, admin.Email)
	err := row.Scan(&storedPHC)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}

	slog.Debug("Result", "PHC", storedPHC)

	PHCData, err := parsePHC(storedPHC)
	if err != nil {
		return err
	}

	hashedPassword := argon2.IDKey(
		[]byte(admin.Password),
		PHCData.Salt,
		PHCData.Time,
		PHCData.Memory,
		PHCData.Parallelism,
		uint32(len(PHCData.Hash)),
		)

	slog.Debug("Hashes", "Hash", hashedPassword, "Stored", PHCData.Hash)

	if subtle.ConstantTimeCompare(hashedPassword, PHCData.Hash) == 1 {
		return nil
	} else {
		return domain.ErrWrongPassword
	}
}

func (d *DB) AdminByEmail(ctx context.Context, email string) (*domain.Admin, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	row := d.DB.QueryRowContext(ctx, getAdminQuery, email)

	var admin domain.Admin
	if err := row.Scan(&admin.User.ID, &admin.User.Name, &admin.User.Email, &admin.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &admin, nil
}
