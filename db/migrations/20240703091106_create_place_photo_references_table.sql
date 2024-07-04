-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS place_photo_references
(
    id         CHAR(36) PRIMARY KEY,
    place_id   CHAR(36) NOT NULL,
    user_id    CHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (place_id) REFERENCES places (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE place_photos
    ADD COLUMN place_photo_reference_id CHAR(36) DEFAULT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE place_photos
    ADD CONSTRAINT fk_place_photos_place_photo_reference_id FOREIGN KEY (place_photo_reference_id) REFERENCES place_photo_references (id);
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
ALTER TABLE place_photos
    DROP FOREIGN KEY fk_place_photos_place_photo_reference_id;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE place_photos
    DROP COLUMN place_photo_reference_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS place_photo_references;
-- +goose StatementEnd


