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
VALUES ($1, $2, $3, $4, 0)
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