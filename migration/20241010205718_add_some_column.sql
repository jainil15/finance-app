-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE users (
  id UUID PRIMARY KEY,
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL unique,
  password_hash varchar(255) NOT NULL
);
-- +goose Down
DROP TABLE users;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
