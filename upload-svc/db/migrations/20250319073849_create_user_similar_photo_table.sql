-- +goose Up
-- +goose StatementBegin
CREATE TYPE similarity_level AS ENUM (
    '1','2','3','4','5','6'
);

CREATE TABLE IF NOT EXISTS user_similar_photos (
    id CHAR(26) PRIMARY KEY NOT NULL,
    photo_id CHAR(26) NOT NULL,
    user_id CHAR(26) NOT NULL,
    similarity similarity_level NOT NULL,
    is_wishlist boolean DEFAULT false,
    is_resend boolean DEFAULT false,
    is_cart boolean DEFAULT false,
    is_favorite boolean DEFAULT false,  
    created_at TIMESTAMPTZ not null,
    updated_at TIMESTAMPTZ not null,
    FOREIGN KEY(photo_id) REFERENCES photos(id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_similar_photos;
DROP TYPE similarity_level;
-- +goose StatementEnd