-- +goose Up
-- +goose StatementBegin
ALTER TABLE plan_candidate_sets
    ADD COLUMN is_place_searched BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE plan_candidate_sets
    DROP COLUMN is_place_searched;
-- +goose StatementEnd
