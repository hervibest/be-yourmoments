-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS facecams (
    id CHAR(26) PRIMARY KEY NOT NULL,
    user_id CHAR(26) NOT NULL,
    file_name VARCHAR(100) NOT NULL,
    file_key VARCHAR(100) NOT NULL,
    title VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    checksum VARCHAR(64),
    url TEXT,
    is_processed BOOLEAN NOT NULL DEFAULT false,
    original_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS facecams;

-- +goose StatementEnd