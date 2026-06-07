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
INSERT INTO users (uuid, name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING id, uuid, name, email, password;

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