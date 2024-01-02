-- +goose Up
-- plan_candidate_sets テーブル
CREATE TABLE plan_candidate_sets
(
    id                            CHAR(36) PRIMARY KEY,
    latitude_start                FLOAT,
    longitude_start               FLOAT,
    created_from_current_location BOOL,
    created_at                    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at                    TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- plan_candidates テーブル
CREATE TABLE plan_candidates
(
    id                    CHAR(36) PRIMARY KEY,
    name                  VARCHAR(255),
    plan_candidate_set_id CHAR(36),
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id)
);

-- plan_candidate_set_searched_places テーブル
CREATE TABLE plan_candidate_set_searched_places
(
    id                    CHAR(36) PRIMARY KEY,
    plan_candidate_set_id CHAR(36),
    place_id              CHAR(36),
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id),
    FOREIGN KEY (place_id) REFERENCES places (id)
);

-- plan_candidate_places テーブル
CREATE TABLE plan_candidate_places
(
    id                CHAR(36) PRIMARY KEY,
    plan_candidate_id CHAR(36),
    place_id          CHAR(36),
    `order`           INT,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_id) REFERENCES plan_candidates (id),
    FOREIGN KEY (place_id) REFERENCES places (id)
);

-- plan_candidate_set_categories テーブル
CREATE TABLE plan_candidate_set_categories
(
    id                    CHAR(36) PRIMARY KEY,
    plan_candidate_set_id CHAR(36),
    category              VARCHAR(255),
    is_selected           BOOL,
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_candidate_set_id) REFERENCES plan_candidate_sets (id)
);

-- +goose Down
DROP TABLE IF EXISTS plan_candidate_places;
DROP TABLE IF EXISTS plan_candidate_set_searched_places;
DROP TABLE IF EXISTS plan_candidate_set_categories;
DROP TABLE IF EXISTS plan_candidates;
DROP TABLE IF EXISTS plan_candidate_sets;