-- +goose Up
-- +goose StatementBegin
ALTER TABLE plans
    ADD COLUMN location POINT NOT NULL DEFAULT (ST_PointFromText('POINT(0 0)'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE plans
    DROP COLUMN location;
-- +goose StatementEnd
