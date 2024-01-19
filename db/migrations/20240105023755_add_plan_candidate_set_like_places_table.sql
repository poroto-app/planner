-- +goose Up
-- +goose StatementBegin
CREATE TABLE plan_candidate_set_like_places
(
    id                    CHAR(36)  NOT NULL PRIMARY KEY,
    plan_candidate_set_id CHAR(36)  NOT NULL,
    place_id              CHAR(36)  NOT NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id),
    FOREIGN KEY (place_id) REFERENCES places (id),
    UNIQUE (plan_candidate_set_id, place_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE plan_candidate_set_like_places;
-- +goose StatementEnd
