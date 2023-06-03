-- name: CreateAccount :one
INSERT INTO
    accounts (
        first_name,
        last_name,
        gender,
        balance,
        currency
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1
FOR NO KEY UPDATE;

-- name: UpdateAccount :one
UPDATE accounts SET balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id
LIMIT $1
OFFSET $2;