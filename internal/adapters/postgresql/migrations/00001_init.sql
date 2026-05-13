-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE, 
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL
);

CREATE INDEX idx_users_uuid ON users(uuid);

CREATE TABLE users_action (
    id BIGSERIAL PRIMARY KEY, 
    users_id BIGINT,
    action TEXT NOT NULL,
    type INTEGER NOT NULL,
    ip TEXT NOT NULL,
    country TEXT NOT NULL,
    action_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_users_action FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE users_log (
    id BIGSERIAL PRIMARY KEY, 
    users_id BIGINT, 
    old_value JSONB,
    new_value JSONB,
    table_name TEXT,
    type INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    log_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_users_log FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE category (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE, 
    name TEXT NOT NULL,
    users_id BIGINT,
    CONSTRAINT fk_category_user FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE expense (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE, 
    title TEXT NOT NULL,
    amount INTEGER NOT NULL, 
    date TIMESTAMPTZ NOT NULL,
    category_id BIGINT,
    users_id BIGINT,
    CONSTRAINT fk_expense_category FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE SET NULL,
    CONSTRAINT fk_expense_user FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE SET NULL
);
-- +goose Down
DROP TABLE IF EXISTS expense;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS users_log;
DROP TABLE IF EXISTS users_action;
DROP TABLE IF EXISTS users;
