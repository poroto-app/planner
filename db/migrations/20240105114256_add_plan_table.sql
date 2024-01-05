-- +goose Up
-- plans テーブル
CREATE TABLE IF NOT EXISTS plans
(
    id         char(36) PRIMARY KEY NOT NULL,
    user_id    char(36),
    name       VARCHAR(2000)        NOT NULL,
    created_at DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME             NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- plan_places テーブル
CREATE TABLE plan_places
(
    id         CHAR(36) PRIMARY KEY,
    plan_id    CHAR(36) NOT NULL,
    place_id   CHAR(36) NOT NULL,
    sort_order INT      NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_id) REFERENCES plans (id),
    FOREIGN KEY (place_id) REFERENCES places (id)
);

-- +goose Down
DROP TABLE IF EXISTS plan_places;

DROP TABLE IF EXISTS plans;
