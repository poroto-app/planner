-- +goose Up
CREATE TABLE plan_collages
(
    id         CHAR(36) PRIMARY KEY,
    plan_id    CHAR(36)  NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_id) REFERENCES plans (id)
);

CREATE TABLE plan_collage_photos
(
    id              CHAR(36) PRIMARY KEY,
    plan_collage_id CHAR(36)  NOT NULL,
    place_id        CHAR(36)  NOT NULL,
    place_photo_id  CHAR(36)  NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_collage_id) REFERENCES plan_collages (id),
    FOREIGN KEY (place_id) REFERENCES places (id),
    FOREIGN KEY (place_photo_id) REFERENCES place_photos (id)
);

-- +goose Down
DROP TABLE plan_collage_photos;
DROP TABLE plan_collages;
