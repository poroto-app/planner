-- +goose Up
START TRANSACTION;

ALTER TABLE plan_candidates
    ADD COLUMN parent_plan_id CHAR(36) DEFAULT NULL;

ALTER TABLE plan_candidates
    ADD CONSTRAINT fk_plan_candidates_parent_plan_id FOREIGN KEY (parent_plan_id) REFERENCES plans(id);

COMMIT;

-- +goose Down
START TRANSACTION;

ALTER TABLE plan_candidates
    DROP FOREIGN KEY fk_plan_candidates_parent_plan_id;

ALTER TABLE plan_candidates
    DROP COLUMN parent_plan_id;

COMMIT;