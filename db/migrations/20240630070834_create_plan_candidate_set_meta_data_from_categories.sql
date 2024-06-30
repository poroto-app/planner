-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plan_candidate_set_meta_data_from_categories
(
    id                    CHAR(36) PRIMARY KEY,
    plan_candidate_set_id CHAR(36)     NOT NULL,
    category_id           VARCHAR(256) NOT NULL,
    range_in_meters       INT          NOT NULL,
    latitude              DOUBLE       NOT NULL,
    longitude             DOUBLE       NOT NULL,
    created_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id),
    INDEX (plan_candidate_set_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS plan_candidate_set_meta_data_from_categories;
-- +goose StatementEnd
