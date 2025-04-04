-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS social_medias (
    id CHAR(26) PRIMARY KEY NOT NULL,
    name VARCHAR(100) NOT NULL UNIQUE,
    base_url TEXT ,
    logo_url TEXT ,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS social_medias;
-- +goose StatementEnd