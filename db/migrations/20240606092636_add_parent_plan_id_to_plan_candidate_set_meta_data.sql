-- +goose Up
-- +goose StatementBegin
ALTER TABLE plan_candidate_set_meta_data
    ADD COLUMN parent_plan_id CHAR(36) DEFAULT NULL,
    ADD CONSTRAINT fk_plan_candidate_set_meta_data_parent_plan_id FOREIGN KEY (parent_plan_id) REFERENCES plans(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE plan_candidate_set_meta_data
    DROP CONSTRAINT fk_plan_candidate_set_meta_data_parent_plan_id,
    DROP COLUMN parent_plan_id;
-- +goose StatementEnd
