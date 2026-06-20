-- +goose Up
create table user_session (
    id uuid primary key,
    user_id integer not null references "users"(id) on delete cascade,
    refresh_token text not null unique,
    created_at timestamp with time zone default now(),
    expires_at timestamp with time zone not null,
    is_active boolean default true
);

-- +goose Down
drop table user_session;
