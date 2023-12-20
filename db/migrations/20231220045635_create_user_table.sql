-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE users
(
    id           VARCHAR(36) PRIMARY KEY,
    firebase_uid VARCHAR(255) UNIQUE NOT NULL,
    name         VARCHAR(255),
    email        VARCHAR(255),
    photo_url    VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE users;
-- +goose StatementEnd
