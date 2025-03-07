-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS photo_profiles (
    id UUID primary key not null,
    user_id UUID UNIQUE,
    size INT,
    url TEXT,
    created_at TIMESTAMPTZ not null,
    updated_at TIMESTAMPTZ not null
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_profiles;
-- +goose StatementEnd