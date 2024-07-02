-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS place_photo_versions
(
    id              CHAR(36)    PRIMARY KEY,
    place_photo_id  CHAR(36)    NOT NULL,
    photo_url             CHAR(36)    NOT NULL,
    width           INT         NOT NULL,
    height          INT         NOT NULL,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (place_photo_id) REFERENCES place_photos (id),
    UNIQUE (place_photo_id, photo_url),
    UNIQUE (place_photo_id, width, height)
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE place_photos
    DROP COLUMN photo_url,
    DROP COLUMN width,
    DROP COLUMN height;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS place_photo_versions;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE place_photos
    ADD COLUMN photo_url CHAR(36) DEFAULT NULL,
    ADD COLUMN width INT DEFAULT 0,
    ADD COLUMN height INT DEFAULT 0;
-- +goose StatementEnd