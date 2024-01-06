-- +goose Up
-- +goose StatementBegin
ALTER TABLE plans
    ADD COLUMN location POINT NOT NULL DEFAULT (ST_PointFromText('POINT(0 0)')),
    ADD SPATIAL INDEX location (location);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE plans
    DROP COLUMN location,
    DROP INDEX location;
-- +goose StatementEnd
