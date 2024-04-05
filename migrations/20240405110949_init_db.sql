-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TYPE token_type_enum AS ENUM ('user', 'admin');

CREATE TABLE IF NOT EXISTS users (
    id              TEXT                     NOT NULL,
    username        TEXT UNIQUE              NOT NULL,
    hashed_password TEXT                     NOT NULL,

    CONSTRAINT pk_users PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tokens (
    id         TEXT                     NOT NULL,
    user_id    TEXT                     NOT NULL,
    value      TEXT                     NOT NULL,
    type token_type_enum default 'user',

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expiration_date  TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT pk_tokens PRIMARY KEY (id),
    CONSTRAINT fk_tokens_users FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS features (
    id int NOT NULL,
    name TEXT,
    CONSTRAINT pk_features PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS banners (
    ID SERIAL NOT NULL,
    feature_id INT,
    content JSONB,
    is_active BOOLEAN,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT pk_banners PRIMARY KEY (id),
    CONSTRAINT fk_banners_features FOREIGN KEY (feature_id) REFERENCES features (id)
);

CREATE TABLE IF NOT EXISTS tags (
    id INT NOT NULL,
    name TEXT,
    CONSTRAINT pk_tags PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS banner_tags (
  banner_id INT REFERENCES banners(id),
  tag_id INT REFERENCES tags (id),
  PRIMARY KEY (banner_id, tag_id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE tokens;
DROP TABLE users;
DROP TABLE banner_tags;
DROP TABLE banners;
DROP TABLE tags;
DROP TABLE features;
DROP TYPE token_type_enum;
