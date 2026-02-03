package db

import (
	"database/sql"
	"inventory/internal/config"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	databaseSchema = `
	CREATE TYPE eboard_position AS ENUM (
	    'President',
	    'VP of Facilities',
	    'Gear Manager',
	    'VP of Activities',
	    'Treasurer',
	    'PR',
	    'Secretary'
	);

	CREATE TABLE IF NOT EXISTS public.users (
		id UUID PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL,
		CONSTRAINT users_email_first_name_last_name_key unique(email, first_name, last_name)
	);

	CREATE TABLE IF NOT EXISTS admins (
		user_id UUID REFERENCES users(id) PRIMARY KEY,
		role eboard_position NOT NULL,
		password_hash VARCHAR(255),
		active BOOLEAN NOT NULL DEFAULT TRUE
	);

	CREATE TABLE IF NOT EXISTS checkouts (
		id SERIAL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id),
		checkout_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		personal BOOLEAN NOT NULL DEFAULT FALSE,
		notes TEXT DEFAULT NULL,
		created_by UUID NOT NULL REFERENCES admins(user_id)
	);

	CREATE TABLE IF NOT EXISTS item_types (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		parent_id INTEGER DEFAULT NULL
	);

	CREATE TABLE IF NOT EXISTS items (
		id VARCHAR(8) PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		item_type_id INTEGER NOT NULL REFERENCES item_types(id),
		notes TEXT DEFAULT NULL,
		purchase_date TIMESTAMP DEFAULT NULL,
		deleted BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS checkout_items (
		checkout_id INTEGER NOT NULL REFERENCES checkouts(id),
		item_id VARCHAR(8) NOT NULL REFERENCES items(id),
		return_date TIMESTAMP DEFAULT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_checkout_items ON checkout_items (checkout_id, item_id);

	CREATE TABLE IF NOT EXISTS schema_version (
		version INTEGER NOT NULL DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS session (
		session_id VARCHAR(64) PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id),
		expires_at TIMESTAMP NOT NULL
	);

	INSERT INTO schema_version (version) VALUES ($1);
	`

	getSchemaVersion = `SELECT version FROM schema_version`
	updateSchemaVersion = `UPDATE schema_version SET version = $1`

	schema_version = 2
)

var (
	migrations = []string{
		`ALTER TABLE items ADD COLUMN date_purchased TIMESTAMP DEFAULT NULL;`,
		`ALTER TABLE item_types ADD UNIQUE (name)`,
	}
)

type DB struct {
	DB *sql.DB
}

func runMigrations(db *sql.DB, start int) (*DB, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	for i := start; i < len(migrations); i++ {
		_, err := tx.Exec(migrations[i])
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	_, err = tx.Exec(updateSchemaVersion, schema_version)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
}

func NewDBWithSchema(conf config.PGConfig) (*DB, error) {
	log.Info("Connecting to database")
	postgresDB, err := sql.Open("postgres", conf.ConnStr())
	if err != nil {
		return nil, err
	}

	var version int
	row := postgresDB.QueryRow(getSchemaVersion)
	if row.Err() == nil {
		row.Scan(&version)
		if version == schema_version {
			return &DB{ DB: postgresDB }, nil
		} else {
			return runMigrations(postgresDB, version)
		}
	}
	_, err = postgresDB.Exec(databaseSchema, schema_version)
	if err != nil {
		return nil, err
	}

	return &DB{ DB: postgresDB }, nil
}
