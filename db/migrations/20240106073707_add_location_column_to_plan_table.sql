-- +goose Up
-- +goose StatementBegin
ALTER TABLE plans
    ADD COLUMN location POINT NOT NULL DEFAULT (ST_PointFromText('POINT(0 0)')),
    ADD SPATIAL INDEX plan_location_index (location);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE plans
    DROP COLUMN location,
    DROP INDEX plan_location_index;
-- +goose StatementEnd
