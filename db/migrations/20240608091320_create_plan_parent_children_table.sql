-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plan_parent_children
(
    id             CHAR(36) PRIMARY KEY,
    parent_plan_id CHAR(36)  NOT NULL,
    child_plan_id  CHAR(36)  NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_plan_id) REFERENCES plans (id),
    FOREIGN KEY (child_plan_id) REFERENCES plans (id),
    UNIQUE (parent_plan_id, child_plan_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS plan_parent_children;
-- +goose StatementEnd
