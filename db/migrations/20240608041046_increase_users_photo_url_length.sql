-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    MODIFY photo_url VARCHAR(2048);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    MODIFY photo_url VARCHAR(255);
-- +goose StatementEnd
