-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS place_recommendations
(
    id         CHAR(36) PRIMARY KEY,
    place_id   CHAR(36) NOT NULL,
    sort_order INT      NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (place_id) REFERENCES places (id),
    UNIQUE KEY place_id (place_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE place_recommendations;
-- +goose StatementEnd
