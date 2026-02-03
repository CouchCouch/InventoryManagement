CREATE TYPE IF NOT EXISTS eboard_position AS ENUM (
    'President',
    'VP of Facilities',
    'Gear Manager',
    'VP of Activities',
    'Treasurer',
    'PR',
    'Secretary'
);

CREATE TYPE IF NOT EXISTS item_type AS ENUM (
    'canoe',
    'sleeping bag',
    'tent',
    'backpack',
    'backpacking stove',
    'paddle',
    'life jacket',
)

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
);

CREATE TABLE IF NOT EXISTS admins (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    active BOOLEAN NOT NULL,
);

CREATE TABLE IF NOT EXISTS checkouts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    checkout_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    personal BOOLEAN NOT NULL,
    trip VARCHAR(100),
    additional_info TEXT,
);

CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY,
    name VARCHAR(100) DEFAULT NULL,
    identifiers TEXT,
    item_type_id INTEGER NOT NULL REFERENCES item_types(id),
    date_purchased DATE,
);

CREATE TABLE IF NOT EXISTS checkout_items (
    checkout_id INTEGER NOT NULL REFERENCES checkouts(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    returned TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (checkout_id, item_id)
);
