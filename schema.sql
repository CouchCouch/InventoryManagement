CREATE TABLE items  (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    quantity INTEGER NOT NULL
);

CREATE TABLE checkouts (
    id SERIAL PRIMARY KEY NOT NULL,
    item_id INTEGER,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    checkout_date DATE NOT NULL,
    returned BOOLEAN DEFAULT FALSE,
    emailed BOOLEAN DEFAULT FALSE,
    email_count INTEGER DEFAULT 0,
    FOREIGN KEY(item_id) REFERENCES items(id)
);
