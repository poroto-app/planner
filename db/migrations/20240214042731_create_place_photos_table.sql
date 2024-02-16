-- +goose Up
-- +goose StatementBegin
CREATE TABLE place_photos (
  id            CHAR(36) PRIMARY KEY,
  place_id      CHAR(36) NOT NULL,
  user_id       VARCHAR(36) NOT NULL,
  photo_url     VARCHAR(255) NOT NULL,
  width         INT NOT NULL,
  height        INT NOT NULL,
  created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (place_id) REFERENCES places (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS place_photos;
-- +goose StatementEnd
