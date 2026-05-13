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