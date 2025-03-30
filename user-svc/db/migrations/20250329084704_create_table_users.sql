-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id CHAR(26) PRIMARY KEY NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    email_verified_at TIMESTAMPTZ,
    password VARCHAR(60),
    phone_number VARCHAR(15) UNIQUE,
    phone_number_verified_at TIMESTAMPTZ,
    google_id VARCHAR(21) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- +goose StatementEnd