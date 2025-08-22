package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	databaseSchema = `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		email VARCHAR(100) NOT NULL,
		deleted BOOLEAN DEFAULT FALSE,
		unique(email, first_name, last_name)
	);

	CREATE TABLE IF NOT EXISTS admins (
		user_id INTEGER REFERENCES users(id),
		role VARCHAR(50) NOT NULL,
		password_hash VARCHAR(255)
	);

	CREATE TABLE IF NOT EXISTS item_types (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS checkouts (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
		checkout_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		notes TEXT DEFAULT NULL,
		created_by INTEGER NOT NULL REFERENCES admins(user_id)
	);

	CREATE TABLE IF NOT EXISTS items (
		id VARCHAR(8) PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		identifiers TEXT,
		item_type_id INTEGER NOT NULL REFERENCES item_types(id),
		notes TEXT DEFAULT NULL,
		deleted BOOLEAN DEFAULT FALSE
	);

	CREATE TABLE IF NOT EXISTS checkout_items (
		checkout_id INTEGER NOT NULL REFERENCES checkouts(id),
		item_id VARCHAR(8) NOT NULL REFERENCES items(id),
		return_date TIMESTAMP DEFAULT NULL,
	);
	CREATE INDEX IF NOT EXISTS idx_checkout_items ON checkout_items (checkout_id, item_id)

	CREATE TABLE IF NOT EXISTS schema_version (
		version INTEGER DEFAULT 0 NOT NULL
	);

	CREATE TABLE IF NOT EXISTS session (
		session_id VARCHAR(64) PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id),
	);

	INSERT INTO admins (role) VALUES
		('President'),
		('Vice President of Activities'),
		('Vice President of Facilities'),
		('Gear Manager'),
		('Treasurer'),
		('Secretary');
	`
)

type db struct {
	DB *sql.DB
}

func NewDBWithSchema(connStr string) (*db, error) {
	log.Info("Connecting to database")
	postgresDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	postgresDB.Exec(databaseSchema)

	db := &db{
		DB: postgresDB,
	}

	return db, nil
}
