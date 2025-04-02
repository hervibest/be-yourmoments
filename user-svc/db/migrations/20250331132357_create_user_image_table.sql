-- +goose Up
-- +goose StatementBegin
CREATE TYPE image_type AS ENUM (
    'PROFILE','COVER'
);

CREATE TABLE IF NOT EXISTS user_images (
    id CHAR(26) PRIMARY KEY NOT NULL,
    user_profile_id CHAR(26) NOT NULL,
    file_name VARCHAR(100) NOT NULL,
    file_key VARCHAR(100) NOT NULL,
    image_type image_type NOT NULL,
    size BIGINT NOT NULL,
    checksum VARCHAR(64),
    url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (user_profile_id) REFERENCES user_profiles(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_images;
DROP TYPE image_type;
-- +goose StatementEnd