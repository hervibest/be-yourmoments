-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS photos (
    id UUID primary key not null,
    user_id UUID UNIQUE,
    owned_by_user_id UUID UNIQUE,
    size INT,
    url TEXT,
    preview_url TEXT,
    preview_with_bounding_url TEXT,
    created_at TIMESTAMPTZ not null,
    updated_at TIMESTAMPTZ not null
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS photos;
-- +goose StatementEnd