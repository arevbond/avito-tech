-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS users (
 id              TEXT                     NOT NULL,
 username        TEXT UNIQUE              NOT NULL,
 hashed_password TEXT                     NOT NULL,
 is_admin BOOLEAN,

 CONSTRAINT pk_users PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tokens (
  id         TEXT                     NOT NULL,
  user_id    TEXT                     NOT NULL,
  value      TEXT                     NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  expiration_date  TIMESTAMP WITH TIME ZONE NOT NULL,

  CONSTRAINT pk_tokens PRIMARY KEY (id),
  CONSTRAINT fk_tokens_users FOREIGN KEY (user_id) REFERENCES users (id)
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
    DROP TABLE tokens;
    DROP TABLE users;
