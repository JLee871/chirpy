-- +goose Up
ALTER TABLE users
ADD hashed_password text not null default 'unset';

-- +goose Down
ALTER TABLE users
DROP COLUMN hashed_password;