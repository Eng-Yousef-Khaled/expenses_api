-- +goose Up
alter table expense drop column uuid;
alter table category drop column uuid;
-- +goose Down
alter table expense add column uuid UUID NOT NULL UNIQUE;
alter table category add column uuid UUID NOT NULL UNIQUE;

