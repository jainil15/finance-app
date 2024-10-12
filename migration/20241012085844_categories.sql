-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE categories(
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  name varchar(255) NOT NULL
);
-- +goose Down
DROP TABLE IF EXISTS categories CASCADE;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
