-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS plan_candidate_set_searched_places;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS plan_candidate_set_searched_places
(
    id                    CHAR(36) PRIMARY KEY,
    plan_candidate_set_id CHAR(36)  NOT NULL,
    place_id              CHAR(36)  NOT NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id),
    FOREIGN KEY (place_id) REFERENCES places (id),
    INDEX (plan_candidate_set_id)
);
-- +goose StatementEnd
