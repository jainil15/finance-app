-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE categories
DROP CONSTRAINT categories_pkey;
ALTER TABLE categories
ADD UNIQUE (user_id, name);
ALTER TABLE categories
ADD PRIMARY KEY (id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
