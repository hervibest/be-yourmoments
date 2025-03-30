-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS email_verifications (
    email VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY(email) REFERENCES users(email)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS email_verifications;
-- +goose StatementEnd