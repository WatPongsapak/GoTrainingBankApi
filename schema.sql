CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    first_name TEXT,
    last_name TEXT
);

CREATE TABLE IF NOT EXISTS bankaccounts(
    id serial PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    account_number TEXT,
    name TEXT,
    balance int
);

CREATE TABLE IF NOT EXISTS keys(
    id serial PRIMARY KEY,
    key TEXT
);