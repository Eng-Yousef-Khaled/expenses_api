-- +goose Up
drop table verification_code;
create table user_verification_code (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(6) NOT NULL CHECK (LENGTH(code) = 6),
    users_id BIGINT,
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '15 minutes',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_verification_code_user FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
drop table user_verification_code;
create table verification_code (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL,
    users_id BIGINT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_verification_code_user FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);
