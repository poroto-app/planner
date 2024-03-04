-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_like_places
(
    id         CHAR(36) PRIMARY KEY,
    user_id    VARCHAR(36) NOT NULL,
    place_id   CHAR(36)    NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (place_id) REFERENCES places (id),
    UNIQUE (user_id, place_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_like_places;
-- +goose StatementEnd
