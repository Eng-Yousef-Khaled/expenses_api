-- +goose Up
alter table users rename hashed_password to password;

-- +goose Down
alter table users rename password to hashed_password;
