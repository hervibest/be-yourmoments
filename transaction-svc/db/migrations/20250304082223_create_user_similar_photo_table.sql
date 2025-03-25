-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_similar_photos (
    id UUID primary key not null,
    user_id UUID UNIQUE,
    photo_id UUID UNIQUE,
    similarity DOUBLE PRECISION,
    created_at TIMESTAMPTZ not null,
    updated_at TIMESTAMPTZ not null,
    FOREIGN KEY(photo_id) REFERENCES photos(id)
);
    
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_similar_photos;
-- +goose StatementEnd