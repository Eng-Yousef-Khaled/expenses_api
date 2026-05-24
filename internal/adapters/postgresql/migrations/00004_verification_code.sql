-- +goose Up
create table verification_code (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL,
    users_id BIGINT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_verification_code_user FOREIGN KEY (users_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE verification_code;
