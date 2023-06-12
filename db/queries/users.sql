-- name: CreateUser :one
INSERT INTO users (
    first_name,
    last_name,
    gender,
    email,
    password
) VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id 
LIMIT $1
OFFSET $2;
