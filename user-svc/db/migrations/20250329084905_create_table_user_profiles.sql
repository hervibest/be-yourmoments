-- +goose Up
-- +goose StatementBegin
CREATE TYPE similarity_level AS ENUM (
    '1','2','3','4','5','6'
);

CREATE TABLE IF NOT EXISTS user_profiles (
    id CHAR(26) PRIMARY KEY NOT NULL,
    user_id CHAR(26) UNIQUE NOT NULL ,
    birth_date DATE,
    nickname VARCHAR(100) UNIQUE,
    biography TEXT,
    profile_url TEXT,
    profile_cover_url TEXT,
    similarity similarity_level DEFAULT '3',
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profiles;
DROP TYPE similarity_level;
-- +goose StatementEnd