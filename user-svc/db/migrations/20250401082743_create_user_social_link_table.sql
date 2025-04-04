-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_social_links (
    user_profile_id CHAR(26),
    social_media_id CHAR(26),
    created_at TIMESTAMPTZ NOT NULl,
    updated_at TIMESTAMPTZ NOT NUll,
    PRIMARY KEY (user_profile_id, social_media_id),
    FOREIGN KEY (user_profile_id) REFERENCES user_profiles(id),
    FOREIGN KEY (social_media_id) REFERENCES social_medias(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_social_links;
-- +goose StatementEnd
