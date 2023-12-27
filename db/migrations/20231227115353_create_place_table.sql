-- +goose Up
CREATE TABLE places
(
    id         CHAR(36)     NOT NULL,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE google_places
(
    google_place_id    VARCHAR(255) NOT NULL,
    place_id           CHAR(36)     NOT NULL,
    name               VARCHAR(255) NOT NULL,
    formatted_address  VARCHAR(255),
    vicinity           VARCHAR(255),
    price_level        INT,
    rating             FLOAT,
    user_ratings_total INT,
    latitude           FLOAT,
    longitude          FLOAT,
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (google_place_id),
    FOREIGN KEY (place_id) REFERENCES places (id)
);

CREATE TABLE google_place_types
(
    id              CHAR(36)     NOT NULL,
    google_place_id VARCHAR(255) NOT NULL,
    type            VARCHAR(255) NOT NULL,
    order_num       INT          NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (google_place_id) REFERENCES google_places (google_place_id)
);

CREATE TABLE google_place_photo_references
(
    photo_reference VARCHAR(255) NOT NULL,
    google_place_id VARCHAR(255) NOT NULL,
    width           INT          NOT NULL,
    height          INT          NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (photo_reference),
    FOREIGN KEY (google_place_id) REFERENCES google_places (google_place_id)
);

CREATE TABLE google_place_photo_attributions
(
    id               CHAR(36)     NOT NULL,
    photo_reference  VARCHAR(255) NOT NULL,
    html_attribution TEXT         NOT NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (photo_reference) REFERENCES google_place_photo_references (photo_reference)
);

CREATE TABLE google_place_photos
(
    id              CHAR(36)     NOT NULL,
    photo_reference VARCHAR(255) NOT NULL,
    width           INT          NOT NULL,
    height          INT          NOT NULL,
    url             VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (photo_reference) REFERENCES google_place_photo_references (photo_reference)
);

CREATE TABLE google_place_reviews
(
    id                       CHAR(36)     NOT NULL,
    google_place_id          VARCHAR(255) NOT NULL,
    author_name              VARCHAR(255),
    author_url               VARCHAR(255),
    author_profile_photo_url VARCHAR(255),
    language                 VARCHAR(255),
    rating                   INT,
    text                     TEXT,
    time                     INT,
    created_at               TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at               TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (google_place_id) REFERENCES google_places (google_place_id)
);

-- +goose Down
DROP TABLE IF EXISTS google_place_reviews;

DROP TABLE IF EXISTS google_place_photos;

DROP TABLE IF EXISTS google_place_photo_attributions;

DROP TABLE IF EXISTS google_place_photo_references;

DROP TABLE IF EXISTS google_place_types;

DROP TABLE IF EXISTS google_places;

DROP TABLE IF EXISTS places;
