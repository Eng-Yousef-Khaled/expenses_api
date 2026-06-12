-- +goose Up
alter table users add column is_verification boolean not null default false;

-- +goose Down
alter table users drop column is_verification;
