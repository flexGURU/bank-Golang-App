-- Improved PostgreSQL Schema

-- Users Table
CREATE TABLE users (
    username VARCHAR PRIMARY KEY,
    hashed_password VARCHAR NOT NULL,
    full_name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password_changed_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00+00:00',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Account Table
CREATE TABLE account (
    id SERIAL PRIMARY KEY,
    owner VARCHAR NOT NULL REFERENCES users(username),
    balance INT NOT NULL,
    currency VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT owner_currency_key UNIQUE (owner, currency)
);

-- Entries Table
CREATE TABLE entries (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL REFERENCES account(id),
    amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Transfer Table
CREATE TABLE transfer (
    id SERIAL PRIMARY KEY,
    from_account_id INT NOT NULL REFERENCES account(id),
    to_account_id INT NOT NULL REFERENCES account(id),
    amount INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Sessions Table
CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    username VARCHAR NOT NULL REFERENCES users(username),
    refresh_token VARCHAR NOT NULL,
    user_agent VARCHAR NOT NULL,
    client_ip VARCHAR NOT NULL,
    is_blocked BOOLEAN NOT NULL DEFAULT false,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes
CREATE INDEX idx_account_owner ON account(owner);
CREATE INDEX idx_entries_account_id ON entries(account_id);
CREATE INDEX idx_transfer_from_account ON transfer(from_account_id);
CREATE INDEX idx_transfer_to_account ON transfer(to_account_id);
CREATE INDEX idx_transfer_from_to_account ON transfer(from_account_id, to_account_id);
