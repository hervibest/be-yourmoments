-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS photos (
    id CHAR(26) PRIMARY KEY NOT NULL,
    creator_id CHAR(26) NOT NULL,
    title VARCHAR(100) NOT NULL,
    owned_by_user_id CHAR(26),
    compressed_url TEXT,
    is_this_you_url TEXT,
    your_moments_url TEXT,
    collection_url TEXT,
    price INT NOT NULL,
    price_str VARCHAR(100) NOT NULL,
    original_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS photos;

-- +goose StatementEnd