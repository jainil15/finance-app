-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TYPE transaction_type AS ENUM('expense', 'income');

CREATE TABLE transactions (
  id UUID NOT NULL PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id),
  category_id UUID NOT NULL REFERENCES categories(id),
  currency varchar(100) NOT NULL, 
  value NUMERIC(12,2) NOT NULL,
  type transaction_type NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS transactions CASCADE;
DROP TYPE IF EXISTS transaction_type CASCADE;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
