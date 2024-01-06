-- +goose Up
-- +goose StatementBegin
ALTER TABLE plans
    ADD COLUMN location POINT    NOT NULL,
    ADD COLUMN geohash  CHAR(12) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE plans
    DROP COLUMN location,
    DROP COLUMN geohash;
-- +goose StatementEnd
