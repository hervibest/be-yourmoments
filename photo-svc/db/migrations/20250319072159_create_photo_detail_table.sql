-- +goose Up
-- +goose StatementBegin
CREATE TYPE photo_type AS ENUM (
    'JPG',
    'PNG',
    'JPEG',
    'RAW',
    'DNG',
    'CR2',
    'CR3',
    'NEF',
    'NRW',
    'ARW',
    'SR2',
    'SRF',
    'ORF',
    'RW2',
    'PEF',
    'RAF',
    'TIFF',
    'HEIF',
    'WEBP',
    'AVIF'
);

CREATE TYPE your_moments_type AS ENUM ('COMPRESSED', 'ISYOU', 'YOU', 'COLLECTION');

CREATE TABLE IF NOT EXISTS photo_details (
    id CHAR(26) PRIMARY KEY NOT NULL,
    photo_id CHAR(26) NOT NULL,
    file_name VARCHAR(100) NOT NULL,
    file_key VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    type photo_type NOT NULL,
    checksum VARCHAR(64),
    width SMALLINT NOT NULL,
    height SMALLINT NOT NULL,
    url TEXT,
    your_moments_type your_moments_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    FOREIGN KEY (photo_id) REFERENCES photos(id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS photo_details;

DROP TYPE photo_type;

DROP TYPE your_moments_type;

-- +goose StatementEnd