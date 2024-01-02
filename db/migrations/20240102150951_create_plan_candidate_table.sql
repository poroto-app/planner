-- +goose Up
-- plan_candidate_sets テーブル
CREATE TABLE plan_candidate_sets
(
    id                            CHAR(36) PRIMARY KEY,
    latitude_start                FLOAT     NOT NULL,
    longitude_start               FLOAT     NOT NULL,
    created_from_current_location BOOL      NOT NULL,
    created_at                    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- plan_candidates テーブル
CREATE TABLE plan_candidates
(
    id                    CHAR(36) PRIMARY KEY,
    name                  VARCHAR(255) NOT NULL,
    plan_candidate_set_id CHAR(36)     NOT NULL,
    created_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id)
);

-- plan_candidate_set_searched_places テーブル
CREATE TABLE plan_candidate_set_searched_places
(
    id                    CHAR(36) PRIMARY KEY,
    plan_candidate_set_id CHAR(36)  NOT NULL,
    place_id              CHAR(36)  NOT NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id),
    FOREIGN KEY (place_id) REFERENCES places (id)
);

-- plan_candidate_places テーブル
CREATE TABLE plan_candidate_places
(
    id                    CHAR(36) PRIMARY KEY,
    plan_candidate_set_id CHAR(36)  NOT NULL,
    place_id              CHAR(36)  NOT NULL,
    `order`               INT       NOT NULL,
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id),
    FOREIGN KEY (place_id) REFERENCES places (id)
);

-- plan_candidate_set_categories テーブル
CREATE TABLE plan_candidate_set_categories
(
    id                    CHAR(36) PRIMARY KEY,
    plan_candidate_set_id CHAR(36)     NOT NULL,
    category              VARCHAR(255) NOT NULL,
    is_selected           BOOL         NOT NULL,
    created_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id)
);

-- +goose Down
DROP TABLE IF EXISTS plan_candidate_places;
DROP TABLE IF EXISTS plan_candidate_set_searched_places;
DROP TABLE IF EXISTS plan_candidate_set_categories;
DROP TABLE IF EXISTS plan_candidates;
DROP TABLE IF EXISTS plan_candidate_sets;