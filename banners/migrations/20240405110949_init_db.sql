-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS features (
    id int NOT NULL,
    name TEXT,
    CONSTRAINT pk_features PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS banners (
    id SERIAL NOT NULL,
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
DROP TABLE banner_tags;
DROP TABLE banners;
DROP TABLE tags;
DROP TABLE features;