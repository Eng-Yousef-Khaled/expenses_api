-- name: ListUsers :many
SELECT 
    * 
FROM 
    users
offset $1 limit $2;
-- name: GetUserByID :one
SELECT
    *
FROM
    users
WHERE
    id = $1;

-- name: CreateUser :one
INSERT INTO users (uuid, name, email, password, is_verification)
VALUES ($1, $2, $3, $4, false)
RETURNING id, uuid, name, email, password, is_verification;

-- name: LoginUser :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: CreateVerificationCode :one
INSERT INTO user_verification_code (code, users_id, expires_at)
VALUES ($1, $2, $3)
RETURNING id, code, users_id, expires_at, created_at;

-- name: SetVerificationStatus :exec
UPDATE users
SET is_verification = $1
WHERE id = $2;

-- name: CheckUserVerificationCode :one
SELECT
    *
FROM
    users u
WHERE
    u.id = $1
    AND u.is_verification = false
    AND EXISTS (
        SELECT 1
        FROM user_verification_code
        WHERE users_id = u.id
          AND code = $2
          AND expires_at > NOW()
    );

-- name: CreateUserSession :one
INSERT INTO user_session (id, user_id, refresh_token, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, refresh_token, created_at, expires_at, is_active;

-- name: GetUserSessionByUserId :one
SELECT
    id, user_id, refresh_token, created_at, expires_at, is_active
FROM
    user_session
WHERE
    user_id = $1
    AND is_active = true;
