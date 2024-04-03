-- +goose Up
ALTER TABLE plans
    ADD COLUMN latitude DOUBLE NOT NULL DEFAULT 0;
ALTER TABLE plans
    ADD COLUMN longitude DOUBLE NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE plans
    DROP COLUMN latitude;
ALTER TABLE plans
    DROP COLUMN longitude;
