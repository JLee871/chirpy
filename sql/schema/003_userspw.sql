-- +goose Up
ALTER TABLE users
ADD hashed_password text not null default 'unset';